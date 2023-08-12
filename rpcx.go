package main

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

const (
	contextPackage    = protogen.GoImportPath("context")
	rpcxClientPackage = protogen.GoImportPath("github.com/cctip/cctip-service-client/rpcclient")
)

// generateFile generates a _grpc.pb.go file containing gRPC service definitions.
func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + ".rpcx.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by david client code gen. DO NOT EDIT.")
	g.P("// versions:")
	g.P("// - protoc-gen-rpcx v", version)
	g.P("// - protoc          ", protocVersion(gen))
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	generateFileContent(gen, file, g)
	return g
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

// generateFileContent generates the gRPC service definitions, excluding the package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Services) == 0 {
		return
	}

	g.P("// Reference imports to suppress errors if they are not otherwise used.")
	g.P("var _ = ", contextPackage.Ident("TODO"))
	g.P("var _ = ", rpcxClientPackage.Ident("GetRpcClient"))
	g.P()

	for _, service := range file.Services {
		genService(gen, file, g, service)
	}
}

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service) {
	serviceName := upperFirstLatter(service.GoName)

	g.P("//================== interface skeleton ===================")
	g.P(fmt.Sprintf(`type %sAble interface {`, serviceName))
	g.P(fmt.Sprintf(`// %sAble can be used for interface verification.`, serviceName))
	g.P()
	for _, method := range service.Methods {
		generateAbleCode(g, method)
	}
	g.P(fmt.Sprintf(`}`))

	g.P()
	//g.P("//================== server skeleton ===================")
	//g.P(fmt.Sprintf(`type %[1]sImpl struct {}
	//
	//	// ServeFor%[1]s starts a server only registers one service.
	//	// You can register more services and only start one server.
	//	// It blocks until the application exits.
	//	func ServeFor%[1]s(addr string) error{
	//		s := server.NewServer()
	//		s.RegisterName("%[1]s", new(%[1]sImpl), "")
	//		return s.Serve("tcp", addr)
	//	}
	//`, serviceName))
	//g.P()
	//for _, method := range service.Methods {
	//	generateServerCode(g, service, method)
	//}

	//g.P()
	g.P("//================== client stub ===================")
	g.P(fmt.Sprintf(`// %[1]s is a client wrapped XClient.
		type %[1]sClient struct{
			service string
		}
		// New%[1]sClient wraps a XClient as %[1]sClient.
		// You can pass a shared XClient object created by NewXClientFor%[1]s.
		func New%[1]sClient(service ...string) *%[1]sClient {
			var serviceName string
			if len(service) < 1{
				serviceName = "%[2]s"
            }else{
				serviceName = service[0]
 			}
			return &%[1]sClient{service: serviceName}
		}

	`, serviceName, strings.Replace(string(file.GoPackageName), "_", "/", -1)))
	for _, method := range service.Methods {
		generateClientCode(g, service, method)
	}

	// one client
	//g.P()
	//g.P("//================== oneclient stub ===================")
	//g.P(fmt.Sprintf(`// %[1]sOneClient is a client wrapped oneClient.
	//	type %[1]sOneClient struct{
	//		serviceName string
	//		oneclient *client.OneClient
	//	}
	//
	//	// New%[1]sOneClient wraps a OneClient as %[1]sOneClient.
	//	// You can pass a shared OneClient object created by NewOneClientFor%[1]s.
	//	func New%[1]sOneClient(oneclient *client.OneClient) *%[1]sOneClient {
	//		return &%[1]sOneClient{
	//			serviceName: "%[1]s",
	//			oneclient: oneclient,
	//		}
	//	}
	//
	//	// ======================================================
	//`, serviceName))
	//for _, method := range service.Methods {
	//	generateOneClientCode(g, service, method)
	//}
}

func generateServerCode(g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method) {
	methodName := upperFirstLatter(method.GoName)
	serviceName := upperFirstLatter(service.GoName)
	inType := g.QualifiedGoIdent(method.Input.GoIdent)
	outType := g.QualifiedGoIdent(method.Output.GoIdent)
	g.P(fmt.Sprintf(`// %s is server rpc method as defined
		func (s *%sImpl) %s(ctx context.Context, args *%s, reply *%s) (err error){
			// TODO: add business logics

			// TODO: setting return values
			*reply = %s{}

			return nil
		}
	`, methodName, serviceName, methodName, inType, outType, outType))
}

func generateAbleCode(g *protogen.GeneratedFile, method *protogen.Method) {
	methodName := upperFirstLatter(method.GoName)
	inType := g.QualifiedGoIdent(method.Input.GoIdent)
	outType := g.QualifiedGoIdent(method.Output.GoIdent)
	g.P(fmt.Sprintf(`// %[1]s is server rpc method as defined
		%[1]s(ctx context.Context, args *%[2]s, reply *%[3]s) (err error)
	`, methodName, inType, outType))
}

func generateClientCode(g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method) {
	methodName := upperFirstLatter(method.GoName)
	serviceName := upperFirstLatter(service.GoName)
	inType := g.QualifiedGoIdent(method.Input.GoIdent)
	outType := g.QualifiedGoIdent(method.Output.GoIdent)
	g.P(fmt.Sprintf(`// %s is client rpc method as defined
		func (c *%sClient) %s(ctx context.Context, args *%s)(reply *%s, err error){
			reply = &%s{}
			err = rpcclient.GetRpcClient().Call(ctx, c.service, "%s",args, reply)
			return reply, err
		}
	`, methodName, serviceName, methodName, inType, outType, outType, method.GoName))
}

func generateOneClientCode(g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method) {
	methodName := upperFirstLatter(method.GoName)
	serviceName := upperFirstLatter(service.GoName)
	inType := g.QualifiedGoIdent(method.Input.GoIdent)
	outType := g.QualifiedGoIdent(method.Output.GoIdent)
	g.P(fmt.Sprintf(`// %s is client rpc method as defined
		func (c *%sOneClient) %s(ctx context.Context, args *%s)(reply *%s, err error){
			reply = &%s{}
			err = c.oneclient.Call(ctx,c.serviceName,"%s",args, reply)
			return reply, err
		}
	`, methodName, serviceName, methodName, inType, outType, outType, method.GoName))
}

// upperFirstLatter make the fisrt charater of given string  upper class
func upperFirstLatter(s string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return strings.ToUpper(string(s[0]))
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}
