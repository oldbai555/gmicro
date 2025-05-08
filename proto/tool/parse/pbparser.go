package parse

import (
	"fmt"
	"github.com/emicklei/proto"
	"github.com/pkg/errors"
	utils "gmicro/common"
	"gmicro/pkg/log"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"strings"
	"text/scanner"
)

type RpcNode struct {
	ApiMethod    string
	AuthType     string
	IgnoreSvrRpc bool
	Options      map[string]string
	Rpc          *proto.RPC
}

type MsgFields struct {
	NormalFields []*proto.NormalField
	MapFields    []*proto.MapField
}

const (
	PbContentNodeTypeUndefined = 0
	PbContentNodeTypeRpc       = 1
	PbContentNodeTypeMsg       = 2
	PbContentNodeTypeEnum      = 3
)

type PbContentNode struct {
	Type  int
	Buf   string
	Name  string
	pos   int
	lines []string
}

func (p *PbContentNode) toLines() {
	p.lines = strings.Split(p.Buf, "\n")
}

func (p *PbContentNode) Pack() {
	p.Buf = strings.Join(p.lines, "\n")
}

func (p *PbContentNode) FindCommentCmd(cmd string) int {
	if !strings.HasPrefix(cmd, "@") {
		cmd = "@" + cmd
	}
	for i := 0; i < len(p.lines); i++ {
		v := p.lines[i]
		t := strings.TrimSpace(v)
		if strings.HasPrefix(t, "//") {
			t = strings.TrimSpace(t[2:])
			if strings.HasPrefix(t, cmd) {
				return i
			} else {
				continue
			}
		} else {
			break
		}
	}
	return -1
}

const (
	pbCommentCmdLeading = "\t"
)

func (p *PbContentNode) AddError(err string, comment string) bool {
	errCmdAt := p.FindCommentCmd("@error")
	cmdLine := fmt.Sprintf("%s// @error:", pbCommentCmdLeading)
	errLine := fmt.Sprintf("%s//  %s", pbCommentCmdLeading, err)
	if comment != "" {
		errLine = fmt.Sprintf("%s %s", errLine, comment)
	}
	if errCmdAt < 0 {
		p.lines = append([]string{cmdLine, errLine}, p.lines...)
	} else {
		found := false
		for i := errCmdAt + 1; i < len(p.lines); i++ {
			v := strings.TrimSpace(p.lines[i])
			vLastIndex := strings.LastIndex(v, "//")
			if vLastIndex == -1 {
				break
			}
			v = strings.TrimSpace(v[vLastIndex+2:])
			items := strings.Split(v, " ")
			if len(items) > 0 {
				v = items[0]
				if v == err {
					found = true
					break
				}
			}
		}
		if !found {
			var lines []string
			x := errCmdAt + 1
			lines = append(lines, p.lines[:x]...)
			lines = append(lines, errLine)
			lines = append(lines, p.lines[x:]...)
			p.lines = lines
		} else {
			return false
		}
	}
	p.Pack()
	return true
}
func (p *PbContentNode) RemoveError(errCodes []string) {
	if len(errCodes) == 0 {
		return
	}
	errCodeMap := map[string]bool{}
	errCmdAt := p.FindCommentCmd("@error")
	for _, errCode := range errCodes {
		errCodeMap[errCode] = true
	}
	var removeIndexs []int
	if errCmdAt >= 0 {
		for i := errCmdAt + 1; i < len(p.lines); i++ {
			v := strings.TrimSpace(p.lines[i])
			vLastIndex := strings.LastIndex(v, "//")
			if vLastIndex != 0 {
				break
			}
			v = strings.TrimSpace(v[vLastIndex+2:])
			if strings.HasPrefix(v, "@") {
				break
			}
			items := strings.Split(v, " ")
			if len(items) > 0 {
				v = items[0]
				if !errCodeMap[v] {
					removeIndexs = append(removeIndexs, i)
					delete(errCodeMap, v)
				}
			}
		}
		removeIndexsLen := len(removeIndexs)
		var lines []string
		if removeIndexsLen > 0 {
			removeIndex := removeIndexs[0]
			lines = append(lines, p.lines[:removeIndex]...)
		}
		for i := 1; i < removeIndexsLen-1; i++ {
			removeLastIndex := removeIndexs[i-1]
			removeIndex := removeIndexs[i]
			removeNextIndex := removeIndexs[i+1]
			lines = append(lines, p.lines[removeLastIndex+1:removeIndex]...)
			lines = append(lines, p.lines[removeIndex+1:removeNextIndex]...)
		}
		if removeIndexsLen > 0 {
			removeIndex := removeIndexs[removeIndexsLen-1]
			lines = append(lines, p.lines[removeIndex+1:]...)
			p.lines = lines
		}
		p.Pack()
	}
}

type PbContentNodeSlice []*PbContentNode

func (s PbContentNodeSlice) Len() int           { return len(s) }
func (s PbContentNodeSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s PbContentNodeSlice) Less(i, j int) bool { return s[i].pos < s[j].pos }

type PbContent struct {
	Nodes []*PbContentNode
}

func (p *PbContent) FindRpcNode(name string) *PbContentNode {
	for _, v := range p.Nodes {
		if v.Type == PbContentNodeTypeRpc && v.Name == name {
			return v
		}
	}
	return nil
}

func (p *PbContent) AppendRpcNodeIfNotExisted(name string, buf string) {
	var insertIndex int
	for index, v := range p.Nodes {
		if v.Type == PbContentNodeTypeRpc {
			if v.Name == name {
				return
			}
			insertIndex = index
		}
	}
	if insertIndex == 0 {
		return
	}
	n := &PbContentNode{
		Type: PbContentNodeTypeRpc,
		Name: name,
		Buf:  buf,
	}
	n.toLines()
	p.Nodes = append(p.Nodes[:insertIndex+1], append([]*PbContentNode{n}, p.Nodes[insertIndex+1:]...)...)
}

func (p *PbContent) AppendModelNodeIfNotExisted(name string, buf string) {
	var insertIndex int
	for index, v := range p.Nodes {
		if v.Type == PbContentNodeTypeMsg {
			if v.Name == name {
				return
			}
			if strings.HasPrefix(v.Name, "Model") {
				insertIndex = index
				if strings.HasSuffix(v.Buf, "{\n") || strings.HasSuffix(v.Buf, "{\n\n") {
					insertIndex++
				}
			}
		}
	}
	if insertIndex == 0 {
		return
	}
	n := &PbContentNode{
		Type: PbContentNodeTypeMsg,
		Name: name,
		Buf:  buf + "\n",
	}
	n.toLines()
	p.Nodes = append(p.Nodes[:insertIndex+1], append([]*PbContentNode{n}, p.Nodes[insertIndex+1:]...)...)
}

func (p *PbContent) AppendMsgNodeIfNotExisted(name string, buf string) {
	nodeLen := len(p.Nodes)
	for _, v := range p.Nodes {
		if v.Type == PbContentNodeTypeMsg && v.Name == name {
			return
		}
	}
	n := &PbContentNode{
		Type: PbContentNodeTypeMsg,
		Name: name,
		Buf:  buf,
	}
	n.toLines()
	if nodeLen > 0 {
		node := p.Nodes[nodeLen-1]
		linesLen := len(node.lines)
		if linesLen > 0 {
			if !strings.HasSuffix(node.Buf, "\n") {
				node.Buf += "\n\n"
			} else if !strings.HasSuffix(node.Buf, "\n\n") {
				node.Buf += "\n"
			}
			node.toLines()
		}
	}
	p.Nodes = append(p.Nodes, n)
}

func (p *PbContent) ToBuf() string {
	var items []string
	for _, v := range p.Nodes {
		v.Pack()
		items = append(items, v.Buf)
	}
	return strings.Replace(strings.Join(items, ""), "    ", "\t", -1)
}

type PbContext struct {
	EnumList      []*proto.Enum
	MsgList       []*proto.Message
	RpcList       []*RpcNode
	ServiceName   string
	GoPackageName string
	ImportList    []string
	PackageName   string
	Content       string
	ProtoFilePath string
	addr2Msg      map[string]*proto.Message
	msgMap        map[string]*proto.Message
	enumMap       map[string]*proto.Enum
	RpcMap        map[string]*RpcNode
	GitPath       string
}

func getPos(c *proto.Comment, s *scanner.Position) int {
	if c != nil {
		return c.Position.Offset
	}
	return s.Offset
}

func (p *PbContext) SplitContent() *PbContent {
	c := &PbContent{}
	for _, v := range p.RpcList {
		n := &PbContentNode{
			Type: PbContentNodeTypeRpc,
			Name: v.Rpc.Name,
			pos:  getPos(v.Rpc.Comment, &v.Rpc.Position),
		}
		c.Nodes = append(c.Nodes, n)
	}

	for _, v := range p.MsgList {
		n := &PbContentNode{
			Type: PbContentNodeTypeMsg,
			Name: v.Name,
			pos:  getPos(v.Comment, &v.Position),
		}
		c.Nodes = append(c.Nodes, n)
	}

	for _, v := range p.EnumList {
		n := &PbContentNode{
			Type: PbContentNodeTypeEnum,
			Name: v.Name,
			pos:  getPos(v.Comment, &v.Position),
		}
		c.Nodes = append(c.Nodes, n)
	}

	sort.Sort(PbContentNodeSlice(c.Nodes))

	var nodes, nodesV2 []*PbContentNode
	if len(c.Nodes) > 0 {
		if c.Nodes[0].pos > 0 {
			nodes = append(nodes, &PbContentNode{
				Type: PbContentNodeTypeUndefined,
				Buf:  p.Content[:c.Nodes[0].pos],
				Name: "",
			})
		}

		for i := 0; i < len(c.Nodes)-1; i++ {
			node := c.Nodes[i]
			next := c.Nodes[i+1]
			node.Buf = p.Content[node.pos:next.pos]
			nodes = append(nodes, node)
		}

		last := c.Nodes[len(c.Nodes)-1]
		if last.pos < len(p.Content) {
			nodes = append(nodes, &PbContentNode{
				Type: PbContentNodeTypeUndefined,
				Buf:  p.Content[last.pos:],
				Name: "",
			})
		}

	} else {

		nodes = append(nodes, &PbContentNode{
			Type: PbContentNodeTypeUndefined,
			Buf:  p.Content,
			Name: "",
		})

	}

	// 代码块这里，按分行切，而不要一行切两半
	for i := 0; i < len(nodes)-1; i++ {
		node := nodes[i]
		next := nodes[i+1]
		if node.Buf == "" {
			continue
		}
		if strings.HasSuffix(node.Buf, "\n") {
			continue
		}
		p := strings.LastIndex(node.Buf, "\n")
		if p > 0 {
			p++
			a := node.Buf[:p]
			b := node.Buf[p:]
			node.Buf = a
			next.Buf = b + next.Buf
		}
	}

	for i, node := range nodes {
		var nextType int
		if i+1 < len(nodes) {
			nextType = nodes[i+1].Type
		}
		if node.Type == PbContentNodeTypeRpc && node.Type != nextType {
			node.toLines()
			nodeLineLen := len(node.lines)
			appendNode := &PbContentNode{
				Type: PbContentNodeTypeUndefined,
				Name: "",
			}
			appendNode.lines = append([]string{""}, node.lines[nodeLineLen-3:]...)
			node.lines = node.lines[:nodeLineLen-3]
			node.Pack()
			appendNode.Pack()
			nodesV2 = append(nodesV2, node, appendNode)
		} else {
			nodesV2 = append(nodesV2, node)
		}
	}

	c.Nodes = nodesV2
	for _, v := range c.Nodes {
		v.toLines()
	}
	return c
}

func IsModel(msg *proto.Message) bool {
	return strings.HasPrefix(msg.Name, "Model")
}

func GetMsgFields(msg *proto.Message) *MsgFields {
	var pv ProtoVisitor
	for _, ei := range msg.Elements {
		ei.Accept(&pv)
	}
	return &MsgFields{
		NormalFields: pv.normalFields,
		MapFields:    pv.mapFields,
	}
}

func GetEnumFields(e *proto.Enum) []*proto.EnumField {
	var pv ProtoVisitor
	for _, ei := range e.Elements {
		ei.Accept(&pv)
	}
	return pv.EnumFields
}

func (p *PbContext) buildMap() {
	p.msgMap = map[string]*proto.Message{}
	p.enumMap = map[string]*proto.Enum{}
	p.RpcMap = map[string]*RpcNode{}
	for _, v := range p.MsgList {
		fullName := p.GetMsgFullName(v)
		p.msgMap[fullName] = v
	}
	for _, v := range p.EnumList {
		fullName := p.GetEnumFullName(v)
		p.enumMap[fullName] = v
	}
	for _, v := range p.RpcList {
		p.RpcMap[v.Rpc.Name] = v
	}
}

func (p *PbContext) GetRpc(name string) *RpcNode {
	return p.RpcMap[name]
}

func (p *PbContext) GetMsg(name string) *proto.Message {
	m := p.msgMap[name]
	if m != nil {
		return m
	}
	return nil
}

func (p *PbContext) GetEnum(name string) *proto.Enum {
	m := p.enumMap[name]
	if m != nil {
		return m
	}
	return nil
}

func (p *PbContext) DumpEnum() {
	for k, v := range p.enumMap {
		if v == nil {
			log.Infof("enum %s, val is nil", k)
		} else {
			log.Infof("enum %s", k)
		}
	}
}

func (p *PbContext) AppendErrCode(nameList []string, errcodeBeg int) (int, error) {
	e := p.GetEnum("ErrCode")
	if e == nil {
		err := errors.New("not found ErrCode enum")
		log.Errorf("err is %v", err)
		return 0, err
	}

	// get the next error code
	fields := GetEnumFields(e)
	maxErrCode := 0
	for _, v := range fields {
		if v.Integer > maxErrCode {
			maxErrCode = v.Integer
		}
	}

	if maxErrCode == 0 {
		if errcodeBeg == 0 {
			err := errors.New("not found any err code")
			log.Errorf("err is %v", err)
			return 0, err
		}
		maxErrCode = errcodeBeg
	}

	b := e.Position.Offset
	pos := -1
	for b < len(p.Content) {
		e := b
		for e < len(p.Content) && p.Content[e] != '\n' {
			e++
		}
		line := p.Content[b:e]
		if strings.TrimSpace(line) == "}" {
			pos = b
			break
		}
		b = e + 1
	}

	var l []string
	for _, name := range nameList {
		existed := false
		for _, v := range fields {
			if v.Integer > maxErrCode {
				maxErrCode = v.Integer
			}
			if v.Name == name {
				log.Warnf("%s existed, skip", name)
				existed = true
			}
		}
		if !existed {
			maxErrCode++
			l = append(l, fmt.Sprintf("\t%s = %d;\n", name, maxErrCode))
		}
	}

	if len(l) == 0 {
		return 0, nil
	}

	if pos > 0 {
		left := p.Content[:pos]
		right := p.Content[pos:]
		p.Content = fmt.Sprintf("%s%s%s", left, strings.Join(l, ""), right)
	} else {
		err := errors.New("invalid enum block, missed }")
		log.Errorf("err is %v", err)
		return 0, err
	}

	return len(l), nil
}

var msgTmpl = `
message %s {
}
`

var CurdTmpl = `

message Model{{.ModelName}} {
	uint64 id = 1;
    uint32 created_at = 2;
    uint32 updated_at = 3;
    uint32 deleted_at = 4;
    uint64 creator_id = 5;
}

message Add{{.ModelName}}{{.Sys}}Req {
    Model{{.ModelName}} data = 1 [(validate.rules).message = {required:true}];
}

message Add{{.ModelName}}{{.Sys}}Rsp {
    Model{{.ModelName}} data = 1;
}

message Update{{.ModelName}}{{.Sys}}Req {
    Model{{.ModelName}} data = 1 [(validate.rules).message = {required:true}];
}

message Update{{.ModelName}}{{.Sys}}Rsp {
}

message Del{{.ModelName}}{{.Sys}}ListReq {
	// @ref_to: Get{{.ModelName}}{{.Sys}}ListReq.ListOption
    lb.ListOption list_option = 1 [(validate.rules).message = {required:true}];
}

message Del{{.ModelName}}{{.Sys}}ListRsp {
}

message Get{{.ModelName}}{{.Sys}}Req {
    uint64 id = 1 [(validate.rules).uint64 = {gt:0}];
}

message Get{{.ModelName}}{{.Sys}}Rsp {
    Model{{.ModelName}} data = 1;
}

message Get{{.ModelName}}{{.Sys}}ListReq {
    enum ListOption {
        ListOptionNil = 0;
    }
    lb.ListOption list_option = 1 [(validate.rules).message = {required:true}];
}

message Get{{.ModelName}}{{.Sys}}ListRsp {
    lb.Paginate paginate = 1;
    repeated Model{{.ModelName}} list = 2;
}

`

func (p *PbContext) AppendBlock(block string) {
	var r string
	if !strings.HasSuffix(p.Content, "\n") {
		r = "\n"
	}
	if !strings.HasPrefix(block, "\n") {
		r = r + "\n"
	}
	p.Content = fmt.Sprintf("%s%s%s", p.Content, r, block)
}

func (p *PbContext) AppendEmptyMsgIfNotExisted(name string) {
	if _, ok := p.msgMap[name]; !ok {
		p.AppendBlock(fmt.Sprintf(msgTmpl, name))
	}
}

func (p *PbContext) AppendMsgIfNotExisted(name, block string) {
	if _, ok := p.msgMap[name]; !ok {
		p.AppendBlock(block)
		p.msgMap[name] = &proto.Message{}
	}
}

func (p *PbContext) AppendEnumIfNotExisted(name, block string) {
	if _, ok := p.enumMap[name]; !ok {
		p.AppendBlock(block)
	}
}

func (p *PbContext) InsertRpcBlock(block string) error {
	if len(p.RpcList) > 0 {
		last := p.RpcList[len(p.RpcList)-1]
		b := last.Rpc.Position.Offset
		pos := -1
		for b < len(p.Content) {
			e := b
			for e < len(p.Content) && p.Content[e] != '\n' {
				e++
			}
			line := p.Content[b:e]
			if strings.TrimSpace(line) == "};" || strings.TrimSpace(line) == "}" {
				pos = e
				break
			}
			b = e + 1
		}
		if pos < 0 {
			return errors.New("invalid proto format, not found rpc block ending }")
		}
		left := p.Content[:pos]
		right := p.Content[pos:]
		p.Content = fmt.Sprintf("%s%s%s", left, block, right)
	} else {
		pos := strings.Index(p.Content, "service ")
		if pos < 0 {
			return errors.New("not found `service ` block")
		}
		b := pos
		for b < len(p.Content) {
			e := b
			for e < len(p.Content) && p.Content[e] != '\n' {
				e++
			}
			line := p.Content[b:e]
			if strings.TrimSpace(line) == "}" {
				// pos = e
				break
			}
			b = e + 1
		}
		left := p.Content[:b]
		right := p.Content[b:]
		block = fmt.Sprintf("\t%s\n", strings.TrimSpace(block))
		p.Content = fmt.Sprintf("%s%s%s", left, block, right)
	}
	return nil
}

func (p *PbContext) GetParent(v proto.Visitee) string {
	var nameList []string
	for v != nil {
		m := p.addr2Msg[fmt.Sprintf("%p", v)]
		if m == nil {
			break
		}
		nameList = append(nameList, m.Name)
		vv := &ProtoVisitor{}
		v.Accept(vv)
		if len(vv.msgList) == len(p.MsgList) || len(vv.msgList) != 1 {
			// 到顶层了
			break
		}
		x := vv.msgList[0]
		v = x.Parent
	}

	if len(nameList) == 0 {
		return ""
	} else if len(nameList) == 1 {
		return nameList[0]
	}

	var rev []string
	for i := len(nameList) - 1; i >= 0; i-- {
		rev = append(rev, nameList[i])
	}
	return strings.Join(rev, ".")
}

func isBuiltInType(typ string) bool {
	switch typ {
	case "string", "uint32", "int32", "uint64", "int64", "bool", "bytes", "float", "double":
		return true
	}
	return false
}

func (p *PbContext) GetParentWithType(typ string, v proto.Visitee) string {
	if isBuiltInType(typ) {
		return ""
	}

	// search up
	for v != nil {
		m := p.addr2Msg[fmt.Sprintf("%p", v)]
		if m == nil {
			break
		}
		for _, y := range m.Elements {
			vv := &ProtoVisitor{}
			y.Accept(vv)
			for _, x := range vv.msgList {
				if x.Name == typ {
					goto OUT
				}
			}
		}
		v = m.Parent
	}

OUT:
	var nameList []string
	for v != nil {
		m := p.addr2Msg[fmt.Sprintf("%p", v)]
		if m == nil {
			break
		}
		nameList = append(nameList, m.Name)
		vv := &ProtoVisitor{}
		v.Accept(vv)
		if len(vv.msgList) != 1 {
			break
		}
		x := vv.msgList[0]
		v = x.Parent
	}

	if len(nameList) == 0 {
		return ""
	} else if len(nameList) == 1 {
		return nameList[0]
	}

	var rev []string
	for i := len(nameList) - 1; i >= 0; i-- {
		rev = append(rev, nameList[i])
	}

	return strings.Join(rev, ".")
}

func (p *PbContext) GetEnumFullName(e *proto.Enum) string {
	parent := p.GetParent(e.Parent)
	var fullName string
	if parent != "" {
		fullName = fmt.Sprintf("%s.%s", parent, e.Name)
	} else {
		fullName = e.Name
	}
	return fullName
}

func (p *PbContext) GetMsgFullName(e *proto.Message) string {
	parent := p.GetParent(e.Parent)
	var fullName string
	if parent != "" {
		fullName = fmt.Sprintf("%s.%s", parent, e.Name)
	} else {
		fullName = e.Name
	}
	return fullName
}

type ProtoVisitor struct {
	msgList      []*proto.Message
	EnumFields   []*proto.EnumField
	normalFields []*proto.NormalField
	mapFields    []*proto.MapField
	apiMethod    string
	authType     string
	ignoreSvrRpc bool
	options      map[string]string
}

func (p *ProtoVisitor) VisitMessage(m *proto.Message) {
	p.msgList = append(p.msgList, m)
}

func (p *ProtoVisitor) VisitService(v *proto.Service) {
}

func (p *ProtoVisitor) VisitSyntax(s *proto.Syntax) {
}

func (p *ProtoVisitor) VisitPackage(pkg *proto.Package) {
}

func (p *ProtoVisitor) VisitOptions(o *proto.Option) {
}

func (p *ProtoVisitor) VisitOption(o *proto.Option) {
	if p.options == nil {
		p.options = map[string]string{}
	}
	if strings.Index(o.Name, "lb.ApiMethod") >= 0 {
		p.apiMethod = o.Constant.Source
	}
	if strings.Index(o.Name, "lb.AuthType") >= 0 {
		p.authType = o.Constant.Source
	}
	if strings.Index(o.Name, "lb.IgnoreSvrRpc") >= 0 {
		p.ignoreSvrRpc = true
	}
	p.options[o.Name] = o.Constant.Source
}

func (p *ProtoVisitor) VisitImport(i *proto.Import) {
}

func (p *ProtoVisitor) VisitNormalField(i *proto.NormalField) {
	p.normalFields = append(p.normalFields, i)
}

func (p *ProtoVisitor) VisitEnumField(i *proto.EnumField) {
	p.EnumFields = append(p.EnumFields, i)
}

func (p *ProtoVisitor) VisitEnum(e *proto.Enum) {
}

func (p *ProtoVisitor) VisitComment(e *proto.Comment) {
}

func (p *ProtoVisitor) VisitOneof(o *proto.Oneof) {
}

func (p *ProtoVisitor) VisitOneofField(o *proto.OneOfField) {
}

func (p *ProtoVisitor) VisitReserved(rs *proto.Reserved) {
}

func (p *ProtoVisitor) VisitRPC(rpc *proto.RPC) {
}

func (p *ProtoVisitor) VisitMapField(f *proto.MapField) {
	p.mapFields = append(p.mapFields, f)
}

func (p *ProtoVisitor) VisitGroup(g *proto.Group) {
}

func (p *ProtoVisitor) VisitExtensions(e *proto.Extensions) {
}

func (p *ProtoVisitor) VisitProto(*proto.Proto) {
}

func walkPb(definition *proto.Proto, ctx *PbContext) {
	handleEnum := func(e *proto.Enum) {
		ctx.EnumList = append(ctx.EnumList, e)
	}

	handleMsg := func(p *proto.Message) {
		ctx.MsgList = append(ctx.MsgList, p)
		ctx.addr2Msg[fmt.Sprintf("%p", p)] = p
	}

	handleRpc := func(m *proto.RPC) {
		var apiMethod, authType = "POST", "user"
		var ignoreSvrRpc = false
		options := map[string]string{}
		for _, opt := range m.Elements {
			v := &ProtoVisitor{}
			opt.Accept(v)
			if v.authType != "" {
				authType = v.authType
			}
			if v.apiMethod != "" {
				apiMethod = v.apiMethod
			}
			if v.ignoreSvrRpc {
				ignoreSvrRpc = v.ignoreSvrRpc
			}
			if v.options != nil {
				for key, val := range v.options {
					options[key] = val
				}
			}
		}
		node := &RpcNode{
			Rpc:          m,
			ApiMethod:    apiMethod,
			AuthType:     authType,
			IgnoreSvrRpc: ignoreSvrRpc,
			Options:      options,
		}
		ctx.RpcList = append(ctx.RpcList, node)
	}
	handleService := func(s *proto.Service) {
		ctx.ServiceName = s.Name
	}
	handleImport := func(i *proto.Import) {
		ctx.ImportList = append(ctx.ImportList, i.Filename)
	}
	proto.Walk(
		definition,
		proto.WithEnum(handleEnum),
		proto.WithMessage(handleMsg),
		proto.WithRPC(handleRpc),
		proto.WithService(handleService),
		proto.WithOption(func(o *proto.Option) {
			if o.Name == "go_package" {
				ctx.GoPackageName = o.Constant.Source
			}
		}),
		proto.WithPackage(func(p *proto.Package) {
			ctx.PackageName = p.Name
		}),
		proto.WithImport(handleImport),
	)
}

func ParsePb(protoFile string) (*PbContext, error) {
	t := SearchImportPb(protoFile)
	if t == "" {
		log.Errorf("not found proto file %s", protoFile)
		return nil, errors.New("not found")
	}
	protoFile = t
	reader, err := os.Open(protoFile)
	if err != nil {
		log.Errorf("can not open proto file %s, error is %v", protoFile, err)
		return nil, err
	}
	defer reader.Close()
	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		log.Errorf("proto parse error %v", err)
		return nil, err
	}
	ctx := &PbContext{addr2Msg: make(map[string]*proto.Message)}
	walkPb(definition, ctx)
	reader.Seek(0, io.SeekStart)
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Errorf("read proto file error %v", err)
		return nil, err
	}
	ctx.Content = string(buf)
	ctx.buildMap()
	ctx.ProtoFilePath = protoFile
	return ctx, nil
}

var sep string
var sepEnvPath string

func init() {
	if runtime.GOOS == "windows" {
		sep = `\`
		sepEnvPath = `;`
	} else {
		sep = "/"
		sepEnvPath = ":"
	}
}

var PbIncPaths []string

func SearchImportPb(impPath string) string {
	if utils.FileExists(impPath) {
		return impPath
	}
	for _, incPath := range PbIncPaths {
		p := fmt.Sprintf("%s%s%s", incPath, sep, impPath)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

func (p *PbContext) WriteContentAndReParse() (*PbContext, error) {
	protoFilePath := p.ProtoFilePath
	err := os.WriteFile(protoFilePath, []byte(p.Content), 0666)
	if err != nil {
		log.Errorf("write to %s err %v", protoFilePath, err)
		return nil, err
	}
	newCtx, err := ParsePb(protoFilePath)
	if err != nil {
		log.Errorf("parse %s err %v", protoFilePath, err)
		return nil, err
	}
	return newCtx, nil
}
