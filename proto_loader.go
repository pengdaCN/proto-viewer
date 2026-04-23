package main

import (
	"archive/tar"
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

type ProtoRegistry struct {
	mu          sync.RWMutex
	files       *protoregistry.Files
	types       *dynamicpb.Types
	loadedDescs map[string][]byte
	loadedTypes map[string]string
}

func NewProtoRegistry() *ProtoRegistry {
	return &ProtoRegistry{
		files:       &protoregistry.Files{},
		types:       dynamicpb.NewTypes(nil),
		loadedDescs: make(map[string][]byte),
		loadedTypes: make(map[string]string),
	}
}

func CheckProtoc() error {
	_, err := exec.LookPath("protoc")
	if err != nil {
		return fmt.Errorf("protoc 未安装或不在 PATH 中")
	}
	return nil
}

func (p *ProtoRegistry) LoadProto(protoContent []byte, fileName string) ([]string, error) {
	if err := CheckProtoc(); err != nil {
		return nil, err
	}

	tmpDir, err := os.MkdirTemp("", "proto-*")
	if err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	protoPath := filepath.Join(tmpDir, fileName)
	if err := os.WriteFile(protoPath, protoContent, 0644); err != nil {
		return nil, fmt.Errorf("写入 proto 文件失败: %v", err)
	}

	descPath := filepath.Join(tmpDir, "descriptor.bin")

	cmd := exec.Command("protoc",
		"--descriptor_set_out="+descPath,
		"--include_imports",
		"-I"+tmpDir,
		protoPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("protoc 编译失败: %s", string(output))
	}

	descData, err := os.ReadFile(descPath)
	if err != nil {
		return nil, fmt.Errorf("读取 descriptor 失败: %v", err)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.loadedDescs[fileName] = descData

	return p.rebuildRegistry()
}

func (p *ProtoRegistry) LoadDirectory(reader io.Reader) ([]string, error) {
	if err := CheckProtoc(); err != nil {
		return nil, err
	}

	tmpDir, err := os.MkdirTemp("", "proto-dir-*")
	if err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := extractTar(reader, tmpDir); err != nil {
		return nil, fmt.Errorf("解压 tar 包失败: %v", err)
	}

	protoFiles, err := findProtoFiles(tmpDir)
	if err != nil {
		return nil, fmt.Errorf("扫描 proto 文件失败: %v", err)
	}
	if len(protoFiles) == 0 {
		return nil, fmt.Errorf("无可加载的 proto 文件")
	}

	localFiles := filterLocalProtoFiles(protoFiles)

	g, err := buildDependencyGraph(tmpDir, localFiles)
	if err != nil {
		return nil, fmt.Errorf("构建依赖图失败: %v", err)
	}

	if cycleFile, hasCycle := g.detectCycles(localFiles); hasCycle {
		return nil, fmt.Errorf("检测到循环依赖: %s", cycleFile)
	}

	googleProtoPath := "assets/google-protobuf"
	descPath := filepath.Join(tmpDir, "descriptor.bin")
	args := []string{
		"--descriptor_set_out=" + descPath,
		"--include_imports",
		"-I" + tmpDir,
		"-I" + googleProtoPath,
	}
	args = append(args, protoFiles...)

	cmd := exec.Command("protoc", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("protoc 编译失败: %s", string(output))
	}

	descData, err := os.ReadFile(descPath)
	if err != nil {
		return nil, fmt.Errorf("读取 descriptor 失败: %v", err)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.loadedDescs = make(map[string][]byte)
	p.loadedDescs["directory"] = descData

	return p.rebuildRegistry()
}

func isGoogleProto(path string) bool {
	return strings.HasPrefix(path, "google/")
}

func filterLocalProtoFiles(protoFiles []string) []string {
	var local []string
	for _, f := range protoFiles {
		if !isGoogleProto(f) {
			local = append(local, f)
		}
	}
	return local
}

func extractTar(reader io.Reader, tmpDir string) error {
	tr := tar.NewReader(reader)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if filepath.Ext(header.Name) != ".proto" {
			continue
		}

		cleanName := strings.ReplaceAll(header.Name, "\\", "/")
		cleanName = strings.TrimPrefix(cleanName, "./")
		cleanName = strings.TrimPrefix(cleanName, "/")

		targetPath := filepath.Join(tmpDir, cleanName)
		if header.FileInfo().IsDir() {
			os.MkdirAll(targetPath, 0755)
			continue
		}

		dir := filepath.Dir(targetPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		if _, err := io.Copy(file, tr); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}
	return nil
}

func findProtoFiles(rootDir string) ([]string, error) {
	var protoFiles []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".proto" {
			relPath, err := filepath.Rel(rootDir, path)
			if err != nil {
				return err
			}
			relPath = strings.ReplaceAll(relPath, "\\", "/")
			protoFiles = append(protoFiles, relPath)
		}
		return nil
	})
	return protoFiles, err
}

func findProtoRoot(protoFiles []string) string {
	if len(protoFiles) == 0 {
		return ""
	}
	parts := strings.Split(filepath.ToSlash(protoFiles[0]), "/")
	if len(parts) <= 1 {
		return ""
	}
	return parts[0]
}

type dependencyGraph struct {
	imports    map[string][]string
	importedBy map[string][]string
}

func buildDependencyGraph(rootDir string, protoFiles []string) (*dependencyGraph, error) {
	g := &dependencyGraph{
		imports:    make(map[string][]string),
		importedBy: make(map[string][]string),
	}

	importRegex := regexp.MustCompile(`import\s+(?:public\s+)?[""]([^""]+)[""]`)

	for _, protoFile := range protoFiles {
		fullPath := filepath.Join(rootDir, protoFile)
		file, err := os.Open(fullPath)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			matches := importRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				importPath := matches[1]
				g.imports[protoFile] = append(g.imports[protoFile], importPath)
				g.importedBy[importPath] = append(g.importedBy[importPath], protoFile)
			}
		}
		file.Close()
	}

	return g, nil
}

func (g *dependencyGraph) findRoots(protoFiles []string) []string {
	var roots []string
	for _, protoFile := range protoFiles {
		if len(g.importedBy[protoFile]) == 0 {
			roots = append(roots, protoFile)
		}
	}
	return roots
}

func (g *dependencyGraph) detectCycles(protoFiles []string) (string, bool) {
	visited := make(map[string]int)
	var path []string

	var dfs func(file string) bool
	dfs = func(file string) bool {
		if visited[file] == 1 {
			path = append(path, file)
			return true
		}
		if visited[file] == 2 {
			return false
		}

		visited[file] = 1
		path = append(path, file)

		for _, imp := range g.imports[file] {
			if !g.isLocalImport(imp) {
				continue
			}

			localImp := g.resolveLocalImport(file, imp)
			for _, protoFile := range protoFiles {
				if g.sameFile(protoFile, localImp) || g.sameFile(protoFile, imp) {
					if dfs(protoFile) {
						return true
					}
					break
				}
			}
		}

		visited[file] = 2
		path = path[:len(path)-1]
		return false
	}

	for _, protoFile := range protoFiles {
		if visited[protoFile] == 0 {
			if dfs(protoFile) {
				return protoFile, true
			}
		}
	}
	return "", false
}

func (g *dependencyGraph) isLocalImport(importPath string) bool {
	return !isGoogleProto(importPath)
}

func (g *dependencyGraph) resolveLocalImport(protoFile, importPath string) string {
	protoDir := filepath.Dir(protoFile)
	if protoDir == "." || protoDir == "" {
		return importPath
	}
	return protoDir + "/" + importPath
}

func (g *dependencyGraph) sameFile(a, b string) bool {
	aName := filepath.Base(a)
	bName := filepath.Base(b)
	return aName == bName || a == b
}

func containsProto(files []string, name string) bool {
	for _, f := range files {
		if f == name || filepath.Base(f) == filepath.Base(name) {
			return true
		}
	}
	return false
}

func topoSort(protoFiles []string, g *dependencyGraph) ([]string, error) {
	inDegree := make(map[string]int)
	for _, f := range protoFiles {
		inDegree[f] = 0
	}

	for _, file := range protoFiles {
		for _, imp := range g.imports[file] {
			localImp := filepath.Join(filepath.Dir(file), imp)
			target := localImp
			if !containsProto(protoFiles, target) {
				target = imp
			}
			if containsProto(protoFiles, target) {
				inDegree[target]++
			}
		}
	}

	var queue []string
	for _, f := range protoFiles {
		if inDegree[f] == 0 {
			queue = append(queue, f)
		}
	}
	sort.Strings(queue)

	var result []string
	for len(queue) > 0 {
		file := queue[0]
		queue = queue[1:]
		result = append(result, file)

		for _, imp := range g.imports[file] {
			localImp := filepath.Join(filepath.Dir(file), imp)
			target := localImp
			if !containsProto(protoFiles, target) {
				target = imp
			}
			if containsProto(protoFiles, target) {
				inDegree[target]--
				if inDegree[target] == 0 {
					queue = append(queue, target)
					sort.Strings(queue)
				}
			}
		}
	}

	if len(result) != len(protoFiles) {
		return nil, fmt.Errorf("检测到循环依赖")
	}

	return result, nil
}

func (p *ProtoRegistry) rebuildRegistry() ([]string, error) {
	if len(p.loadedDescs) == 0 {
		p.files = &protoregistry.Files{}
		p.types = dynamicpb.NewTypes(nil)
		p.loadedTypes = make(map[string]string)
		return []string{}, nil
	}

	var fds descriptorpb.FileDescriptorSet
	for _, descData := range p.loadedDescs {
		var fileFds descriptorpb.FileDescriptorSet
		if err := proto.Unmarshal(descData, &fileFds); err != nil {
			return nil, fmt.Errorf("解析 descriptor 失败: %v", err)
		}
		fds.File = append(fds.File, fileFds.File...)
	}

	files, err := protodesc.NewFiles(&fds)
	if err != nil {
		return nil, fmt.Errorf("创建文件描述符失败: %v", err)
	}

	p.files = files
	p.types = dynamicpb.NewTypes(files)

	types := []string{}
	p.loadedTypes = make(map[string]string)

	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Messages().Len(); i++ {
			msgDesc := fd.Messages().Get(i)
			fullName := string(msgDesc.FullName())
			types = append(types, fullName)
			p.loadedTypes[fullName] = ""
		}
		return true
	})

	return types, nil
}

func (p *ProtoRegistry) GetLoadedTypes() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var types []string
	if p.files == nil {
		return types
	}
	p.files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Messages().Len(); i++ {
			msgDesc := fd.Messages().Get(i)
			types = append(types, string(msgDesc.FullName()))
		}
		return true
	})
	sort.Strings(types)
	return types
}

type PaginatedTypes struct {
	Types      []string
	Total      int
	Page       int
	PageSize   int
	TotalPages int
}

func (p *ProtoRegistry) GetLoadedTypesPaginated(page, pageSize int, search string) (*PaginatedTypes, error) {
	allTypes := p.GetLoadedTypes()

	if search != "" {
		lowerSearch := strings.ToLower(search)
		filtered := []string{}
		for _, t := range allTypes {
			if strings.Contains(strings.ToLower(t), lowerSearch) {
				filtered = append(filtered, t)
			}
		}
		allTypes = filtered
	}

	total := len(allTypes)

	if total == 0 {
		return &PaginatedTypes{
			Types:      []string{},
			Total:      0,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: 0,
		}, nil
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	totalPages := (total + pageSize - 1) / pageSize
	if page > totalPages {
		page = totalPages
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}

	return &PaginatedTypes{
		Types:      allTypes[start:end],
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (p *ProtoRegistry) UnloadByName(typeName string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var fileNameToRemove string
	for fname, desc := range p.loadedDescs {
		var fds descriptorpb.FileDescriptorSet
		if err := proto.Unmarshal(desc, &fds); err != nil {
			continue
		}
		for _, fd := range fds.File {
			for _, msg := range fd.MessageType {
				fullName := string(fd.GetPackage()) + "." + msg.GetName()
				if fullName == typeName || fd.GetPackage() == "" && msg.GetName() == typeName {
					fileNameToRemove = fname
					break
				}
			}
			if fileNameToRemove != "" {
				break
			}
		}
		if fileNameToRemove != "" {
			break
		}
	}

	if fileNameToRemove == "" {
		return fmt.Errorf("未找到类型: %s", typeName)
	}

	delete(p.loadedDescs, fileNameToRemove)
	_, err := p.rebuildRegistryLocked()
	return err
}

func (p *ProtoRegistry) rebuildRegistryLocked() ([]string, error) {
	if len(p.loadedDescs) == 0 {
		p.files = &protoregistry.Files{}
		p.types = dynamicpb.NewTypes(nil)
		p.loadedTypes = make(map[string]string)
		return []string{}, nil
	}

	var fds descriptorpb.FileDescriptorSet
	for _, descData := range p.loadedDescs {
		var fileFds descriptorpb.FileDescriptorSet
		if err := proto.Unmarshal(descData, &fileFds); err != nil {
			return nil, fmt.Errorf("解析 descriptor 失败: %v", err)
		}
		fds.File = append(fds.File, fileFds.File...)
	}

	files, err := protodesc.NewFiles(&fds)
	if err != nil {
		return nil, fmt.Errorf("创建文件描述符失败: %v", err)
	}

	p.files = files
	p.types = dynamicpb.NewTypes(files)

	types := []string{}
	p.loadedTypes = make(map[string]string)

	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		for i := 0; i < fd.Messages().Len(); i++ {
			msgDesc := fd.Messages().Get(i)
			fullName := string(msgDesc.FullName())
			types = append(types, fullName)
			p.loadedTypes[fullName] = ""
		}
		return true
	})

	return types, nil
}

func extractGrpcWebPayload(data []byte) ([]byte, error) {
	var result []byte
	offset := 0
	hasDataFrame := false

	for offset < len(data) {
		if offset+5 > len(data) {
			break
		}

		flags := data[offset]
		length := binary.BigEndian.Uint32(data[offset+1 : offset+5])
		payloadStart := offset + 5
		payloadEnd := payloadStart + int(length)

		if payloadEnd > len(data) {
			break
		}

		if flags&0x01 != 0 {
			return nil, errors.New("不支持压缩的 gRPC-Web 数据")
		}

		if flags&0x80 == 0 {
			hasDataFrame = true
			result = append(result, data[payloadStart:payloadEnd]...)
		}

		offset = payloadEnd
	}

	if !hasDataFrame {
		return nil, errors.New("未找到有效的 gRPC-Web DATA 帧")
	}

	return result, nil
}

func (p *ProtoRegistry) Decode(data []byte, typeName string) (string, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	msgType, err := p.types.FindMessageByName(protoreflect.FullName(typeName))
	if err != nil {
		return "", fmt.Errorf("未找到类型 %s: %w", typeName, err)
	}

	dynMsg := dynamicpb.NewMessage(msgType.Descriptor())

	if err := proto.Unmarshal(data, dynMsg); err != nil {
		grpcPayload, grpcErr := extractGrpcWebPayload(data)
		if grpcErr != nil {
			return "", fmt.Errorf("反序列化失败: %v (原始数据) -> %v (gRPC-Web解析: %v)", err, grpcErr, grpcErr)
		}
		dynMsg = dynamicpb.NewMessage(msgType.Descriptor())
		if err2 := proto.Unmarshal(grpcPayload, dynMsg); err2 != nil {
			return "", fmt.Errorf("反序列化失败: %v (原始数据) -> %v (gRPC-Web解析后) -> %v", err, grpcErr, err2)
		}
	}

	jsonBytes, err := protojson.MarshalOptions{
		Indent: "  ",
	}.Marshal(dynMsg)
	if err != nil {
		return "", fmt.Errorf("JSON 转换失败: %v", err)
	}

	return string(jsonBytes), nil
}
