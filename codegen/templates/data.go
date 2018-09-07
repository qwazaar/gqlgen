package templates

var data = map[string]string{
	"args.gotpl":      "\targs := map[string]interface{}{}\n\t{{- range $i, $arg := . }}\n\t\tvar arg{{$i}} {{$arg.Signature }}\n\t\tif tmp, ok := rawArgs[{{$arg.GQLName|quote}}]; ok {\n\t\t\tvar err error\n\t\t\t{{$arg.Unmarshal (print \"arg\" $i) \"tmp\" }}\n\t\t\tif err != nil {\n\t\t\t\treturn nil, err\n\t\t\t}\n\t\t}\n\t\targs[{{$arg.GQLName|quote}}] = arg{{$i}}\n\t{{- end }}\n\treturn args, nil\n",
	"field.gotpl":     "{{ $field := . }}\n{{ $object := $field.Object }}\n\n{{- if $object.Stream }}\n\tfunc (ec *executionContext) _{{$object.GQLType}}_{{$field.GQLName}}(ctx context.Context, field graphql.CollectedField) func() graphql.Marshaler {\n\t\t{{- if $field.Args }}\n\t\t\trawArgs := field.ArgumentMap(ec.Variables)\n\t\t\targs, err := {{ $field.ArgsFunc }}(rawArgs)\n\t\t\tif err != nil {\n\t\t\t\tec.Error(ctx, err)\n\t\t\t\treturn nil\n\t\t\t}\n\t\t{{- end }}\n\t\tctx = graphql.WithResolverContext(ctx, &graphql.ResolverContext{\n\t\t\tField: field,\n\t\t})\n\t\tresults, err := ec.resolvers.{{ $field.ShortInvocation }}\n\t\tif err != nil {\n\t\t\tec.Error(ctx, err)\n\t\t\treturn nil\n\t\t}\n\t\treturn func() graphql.Marshaler {\n\t\t\tres, ok := <-results\n\t\t\tif !ok {\n\t\t\t\treturn nil\n\t\t\t}\n\t\t\tvar out graphql.OrderedMap\n\t\t\tout.Add(field.Alias, func() graphql.Marshaler { {{ $field.WriteJson }} }())\n\t\t\treturn &out\n\t\t}\n\t}\n{{ else }}\n\t// nolint: vetshadow\n\tfunc (ec *executionContext) _{{$object.GQLType}}_{{$field.GQLName}}(ctx context.Context, field graphql.CollectedField, {{if not $object.Root}}obj *{{$object.FullName}}{{end}}) graphql.Marshaler {\n\t\t{{- if $field.Args }}\n\t\t\trawArgs := field.ArgumentMap(ec.Variables)\n\t\t\targs, err := {{ $field.ArgsFunc }}(rawArgs)\n\t\t\tif err != nil {\n\t\t\t\tec.Error(ctx, err)\n\t\t\t\treturn graphql.Null\n\t\t\t}\n\t\t{{- end }}\n\t\trctx := &graphql.ResolverContext{\n\t\t\tObject: {{$object.GQLType|quote}},\n\t\t\tArgs: {{if $field.Args }}args{{else}}nil{{end}},\n\t\t\tField: field,\n\t\t}\n\t\tctx = graphql.WithResolverContext(ctx, rctx)\n\t\tresTmp := ec.FieldMiddleware(ctx, {{if $object.Root}}nil{{else}}obj{{end}}, func(ctx context.Context) (interface{}, error) {\n\t\t\t{{- if $field.IsResolver }}\n\t\t\t\treturn ec.resolvers.{{ $field.ShortInvocation }}\n\t\t\t{{- else if $field.IsMethod }}\n\t\t\t\t{{- if $field.NoErr }}\n\t\t\t\t\treturn {{$field.GoReceiverName}}.{{$field.GoFieldName}}({{ $field.CallArgs }}), nil\n\t\t\t\t{{- else }}\n\t\t\t\t\treturn {{$field.GoReceiverName}}.{{$field.GoFieldName}}({{ $field.CallArgs }})\n\t\t\t\t{{- end }}\n\t\t\t{{- else if $field.IsVariable }}\n\t\t\t\treturn {{$field.GoReceiverName}}.{{$field.GoFieldName}}, nil\n\t\t\t{{- end }}\n\t\t})\n\t\tif resTmp == nil {\n\t\t\t{{- if $field.ASTType.NonNull }}\n\t\t\t\tif !ec.HasError(rctx) {\n\t\t\t\t\tec.Errorf(ctx, \"must not be null\")\n\t\t\t\t}\n\t\t\t{{- end }}\n\t\t\treturn graphql.Null\n\t\t}\n\t\tres := resTmp.({{$field.Signature}})\n\t\trctx.Result = res\n\t\t{{ $field.WriteJson }}\n\t}\n{{ end }}\n",
	"generated.gotpl": "// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.\n\npackage {{ .PackageName }}\n\nimport (\n{{- range $import := .Imports }}\n\t{{- $import.Write }}\n{{ end }}\n)\n\n// NewExecutableSchema creates an ExecutableSchema from the ResolverRoot interface.\nfunc NewExecutableSchema(cfg Config) graphql.ExecutableSchema {\n\treturn &executableSchema{\n\t\tresolvers: cfg.Resolvers,\n\t\tdirectives: cfg.Directives,\n\t\tcomplexity: cfg.Complexity,\n\t}\n}\n\ntype Config struct {\n\tResolvers  ResolverRoot\n\tDirectives DirectiveRoot\n\tComplexity ComplexityRoot\n}\n\ntype ResolverRoot interface {\n{{- range $object := .Objects -}}\n\t{{ if $object.HasResolvers -}}\n\t\t{{$object.GQLType}}() {{$object.GQLType}}Resolver\n\t{{ end }}\n{{- end }}\n}\n\ntype DirectiveRoot struct {\n{{ range $directive := .Directives }}\n\t{{ $directive.Declaration }}\n{{ end }}\n}\n\ntype ComplexityRoot struct {\n{{ range $object := .Objects }}\n\t{{ if not $object.IsReserved -}}\n\t\t{{ $object.GQLType|toCamel }} struct {\n\t\t{{ range $field := $object.Fields -}}\n\t\t\t{{ if not $field.IsReserved -}}\n\t\t\t\t{{ $field.GQLName|toCamel }} {{ $field.ComplexitySignature }}\n\t\t\t{{ end }}\n\t\t{{- end }}\n\t\t}\n\t{{- end }}\n{{ end }}\n}\n\n{{ range $object := .Objects -}}\n\t{{ if $object.HasResolvers }}\n\t\ttype {{$object.GQLType}}Resolver interface {\n\t\t{{ range $field := $object.Fields -}}\n\t\t\t{{ $field.ShortResolverDeclaration }}\n\t\t{{ end }}\n\t\t}\n\t{{- end }}\n{{- end }}\n\n{{ range $object := .Objects -}}\n\t{{ range $field := $object.Fields -}}\n\t\t{{ if $field.Args }}\n\t\t\tfunc {{ $field.ArgsFunc }}(rawArgs map[string]interface{}) (map[string]interface{}, error) {\n\t\t\t{{ template \"args.gotpl\" $field.Args }}\n\t\t\t}\n\t\t{{ end }}\n\t{{ end }}\n{{- end }}\n\n{{ range $directive := .Directives }}\n\t{{ if $directive.Args }}\n\t\tfunc {{ $directive.ArgsFunc }}(rawArgs map[string]interface{}) (map[string]interface{}, error) {\n\t\t{{ template \"args.gotpl\" $directive.Args }}\n\t\t}\n\t{{ end }}\n{{ end }}\n\ntype executableSchema struct {\n\tresolvers  ResolverRoot\n\tdirectives DirectiveRoot\n\tcomplexity ComplexityRoot\n}\n\nfunc (e *executableSchema) Schema() *ast.Schema {\n\treturn parsedSchema\n}\n\nfunc (e *executableSchema) Complexity(typeName, field string, childComplexity int, rawArgs map[string]interface{}) (int, bool) {\n\tswitch typeName + \".\" + field {\n\t{{ range $object := .Objects }}\n\t\t{{ if not $object.IsReserved }}\n\t\t\t{{ range $field := $object.Fields }}\n\t\t\t\t{{ if not $field.IsReserved }}\n\t\t\t\t\tcase \"{{$object.GQLType}}.{{$field.GQLName}}\":\n\t\t\t\t\t\tif e.complexity.{{$object.GQLType|toCamel}}.{{$field.GQLName|toCamel}} == nil {\n\t\t\t\t\t\t\tbreak\n\t\t\t\t\t\t}\n\t\t\t\t\t\t{{ if $field.Args }}\n\t\t\t\t\t\t\targs, err := {{ $field.ArgsFunc }}(rawArgs)\n\t\t\t\t\t\t\tif err != nil {\n\t\t\t\t\t\t\t\treturn 0, false\n\t\t\t\t\t\t\t}\n\t\t\t\t\t\t{{ end }}\n\t\t\t\t\t\treturn e.complexity.{{$object.GQLType|toCamel}}.{{$field.GQLName|toCamel}}(childComplexity{{if $field.Args}}, {{$field.ComplexityArgs}} {{end}}), true\n\t\t\t\t{{ end }}\n\t\t\t{{ end }}\n\t\t{{ end }}\n\t{{ end }}\n\t}\n\treturn 0, false\n}\n\nfunc (e *executableSchema) Query(ctx context.Context, op *ast.OperationDefinition) *graphql.Response {\n\t{{- if .QueryRoot }}\n\t\tec := executionContext{graphql.GetRequestContext(ctx), e}\n\n\t\tbuf := ec.RequestMiddleware(ctx, func(ctx context.Context) []byte {\n\t\t\tdata := ec._{{.QueryRoot.GQLType}}(ctx, op.SelectionSet)\n\t\t\tvar buf bytes.Buffer\n\t\t\tdata.MarshalGQL(&buf)\n\t\t\treturn buf.Bytes()\n\t\t})\n\n\t\treturn &graphql.Response{\n\t\t\tData:       buf,\n\t\t\tErrors:     ec.Errors,\n\t\t\tExtensions: ec.Extensions,\t\t}\n\t{{- else }}\n\t\treturn graphql.ErrorResponse(ctx, \"queries are not supported\")\n\t{{- end }}\n}\n\nfunc (e *executableSchema) Mutation(ctx context.Context, op *ast.OperationDefinition) *graphql.Response {\n\t{{- if .MutationRoot }}\n\t\tec := executionContext{graphql.GetRequestContext(ctx), e}\n\n\t\tbuf := ec.RequestMiddleware(ctx, func(ctx context.Context) []byte {\n\t\t\tdata := ec._{{.MutationRoot.GQLType}}(ctx, op.SelectionSet)\n\t\t\tvar buf bytes.Buffer\n\t\t\tdata.MarshalGQL(&buf)\n\t\t\treturn buf.Bytes()\n\t\t})\n\n\t\treturn &graphql.Response{\n\t\t\tData:       buf,\n\t\t\tErrors:     ec.Errors,\n\t\t\tExtensions: ec.Extensions,\n\t\t}\n\t{{- else }}\n\t\treturn graphql.ErrorResponse(ctx, \"mutations are not supported\")\n\t{{- end }}\n}\n\nfunc (e *executableSchema) Subscription(ctx context.Context, op *ast.OperationDefinition) func() *graphql.Response {\n\t{{- if .SubscriptionRoot }}\n\t\tec := executionContext{graphql.GetRequestContext(ctx), e}\n\n\t\tnext := ec._{{.SubscriptionRoot.GQLType}}(ctx, op.SelectionSet)\n\t\tif ec.Errors != nil {\n\t\t\treturn graphql.OneShot(&graphql.Response{Data: []byte(\"null\"), Errors: ec.Errors})\n\t\t}\n\n\t\tvar buf bytes.Buffer\n\t\treturn func() *graphql.Response {\n\t\t\tbuf := ec.RequestMiddleware(ctx, func(ctx context.Context) []byte {\n\t\t\t\tbuf.Reset()\n\t\t\t\tdata := next()\n\n\t\t\t\tif data == nil {\n\t\t\t\t\treturn nil\n\t\t\t\t}\n\t\t\t\tdata.MarshalGQL(&buf)\n\t\t\t\treturn buf.Bytes()\n\t\t\t})\n\n\t\t\treturn &graphql.Response{\n\t\t\t\tData:       buf,\n\t\t\t\tErrors:     ec.Errors,\n\t\t\t\tExtensions: ec.Extensions,\n\t\t\t}\n\t\t}\n\t{{- else }}\n\t\treturn graphql.OneShot(graphql.ErrorResponse(ctx, \"subscriptions are not supported\"))\n\t{{- end }}\n}\n\ntype executionContext struct {\n\t*graphql.RequestContext\n\t*executableSchema\n}\n\n{{- range $object := .Objects }}\n\t{{ template \"object.gotpl\" $object }}\n\n\t{{- range $field := $object.Fields }}\n\t\t{{ template \"field.gotpl\" $field }}\n\t{{ end }}\n{{- end}}\n\n{{- range $interface := .Interfaces }}\n\t{{ template \"interface.gotpl\" $interface }}\n{{- end }}\n\n{{- range $input := .Inputs }}\n\t{{ template \"input.gotpl\" $input }}\n{{- end }}\n\nfunc (ec *executionContext) FieldMiddleware(ctx context.Context, obj interface{}, next graphql.Resolver) (ret interface{}) {\n\tdefer func() {\n\t\tif r := recover(); r != nil {\n\t\t\tec.Error(ctx, ec.Recover(ctx, r))\n\t\t\tret = nil\n\t\t}\n\t}()\n\t{{- if .Directives }}\n\trctx := graphql.GetResolverContext(ctx)\n\tfor _, d := range rctx.Field.Definition.Directives {\n\t\tswitch d.Name {\n\t\t{{- range $directive := .Directives }}\n\t\tcase \"{{$directive.Name}}\":\n\t\t\tif ec.directives.{{$directive.Name|ucFirst}} != nil {\n\t\t\t\t{{- if $directive.Args }}\n\t\t\t\t\trawArgs := d.ArgumentMap(ec.Variables)\n\t\t\t\t\targs, err := {{ $directive.ArgsFunc }}(rawArgs)\n\t\t\t\t\tif err != nil {\n\t\t\t\t\t\tec.Error(ctx, err)\n\t\t\t\t\t\treturn nil\n\t\t\t\t\t}\n\t\t\t\t{{- end }}\n\t\t\t\tn := next\n\t\t\t\tnext = func(ctx context.Context) (interface{}, error) {\n\t\t\t\t\treturn ec.directives.{{$directive.Name|ucFirst}}({{$directive.CallArgs}})\n\t\t\t\t}\n\t\t\t}\n\t\t{{- end }}\n\t\t}\n\t}\n\t{{- end }}\n\tres, err := ec.ResolverMiddleware(ctx, next)\n\tif err != nil {\n\t\tec.Error(ctx, err)\n\t\treturn nil\n\t}\n\treturn res\n}\n\nfunc (ec *executionContext) introspectSchema() *introspection.Schema {\n\treturn introspection.WrapSchema(parsedSchema)\n}\n\nfunc (ec *executionContext) introspectType(name string) *introspection.Type {\n\treturn introspection.WrapTypeFromDef(parsedSchema, parsedSchema.Types[name])\n}\n\nvar parsedSchema = gqlparser.MustLoadSchema(\n\t&ast.Source{Name: {{.SchemaFilename|quote}}, Input: {{.SchemaRaw|rawQuote}}},\n)\n",
	"input.gotpl":     "\t{{- if .IsMarshaled }}\n\tfunc Unmarshal{{ .GQLType }}(v interface{}) ({{.FullName}}, error) {\n\t\tvar it {{.FullName}}\n\t\tvar asMap = v.(map[string]interface{})\n\t\t{{ range $field := .Fields}}\n\t\t\t{{- if $field.Default}}\n\t\t\t\tif _, present := asMap[{{$field.GQLName|quote}}] ; !present {\n\t\t\t\t\tasMap[{{$field.GQLName|quote}}] = {{ $field.Default | dump }}\n\t\t\t\t}\n\t\t\t{{- end}}\n\t\t{{- end }}\n\n\t\tfor k, v := range asMap {\n\t\t\tswitch k {\n\t\t\t{{- range $field := .Fields }}\n\t\t\tcase {{$field.GQLName|quote}}:\n\t\t\t\tvar err error\n\t\t\t\t{{ $field.Unmarshal (print \"it.\" $field.GoFieldName) \"v\" }}\n\t\t\t\tif err != nil {\n\t\t\t\t\treturn it, err\n\t\t\t\t}\n\t\t\t{{- end }}\n\t\t\t}\n\t\t}\n\n\t\treturn it, nil\n\t}\n\t{{- end }}\n",
	"interface.gotpl": "{{- $interface := . }}\n\nfunc (ec *executionContext) _{{$interface.GQLType}}(ctx context.Context, sel ast.SelectionSet, obj *{{$interface.FullName}}) graphql.Marshaler {\n\tswitch obj := (*obj).(type) {\n\tcase nil:\n\t\treturn graphql.Null\n\t{{- range $implementor := $interface.Implementors }}\n\t\t{{- if $implementor.ValueReceiver }}\n\t\t\tcase {{$implementor.FullName}}:\n\t\t\t\treturn ec._{{$implementor.GQLType}}(ctx, sel, &obj)\n\t\t{{- end}}\n\t\tcase *{{$implementor.FullName}}:\n\t\t\treturn ec._{{$implementor.GQLType}}(ctx, sel, obj)\n\t{{- end }}\n\tdefault:\n\t\tpanic(fmt.Errorf(\"unexpected type %T\", obj))\n\t}\n}\n",
	"models.gotpl":    "// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.\n\npackage {{ .PackageName }}\n\nimport (\n{{- range $import := .Imports }}\n\t{{- $import.Write }}\n{{ end }}\n)\n\n{{ range $model := .Models }}\n\t{{with .Description}} {{.|prefixLines \"// \"}} {{end}}\n\t{{- if .IsInterface }}\n\t\ttype {{.GoType}} interface {}\n\t{{- else }}\n\t\ttype {{.GoType}} struct {\n\t\t\t{{- range $field := .Fields }}\n\t\t\t\t{{- with .Description}}\n\t\t\t\t\t{{.|prefixLines \"// \"}}\n\t\t\t\t{{- end}}\n\t\t\t\t{{- if $field.GoFieldName }}\n\t\t\t\t\t{{ $field.GoFieldName }} {{$field.Signature}} `json:\"{{$field.GQLName}}\"`\n\t\t\t\t{{- else }}\n\t\t\t\t\t{{ $field.GoFKName }} {{$field.GoFKType}}\n\t\t\t\t{{- end }}\n\t\t\t{{- end }}\n\t\t}\n\t{{- end }}\n{{- end}}\n\n{{ range $enum := .Enums }}\n\t{{with .Description}}{{.|prefixLines \"// \"}} {{end}}\n\ttype {{.GoType}} string\n\tconst (\n\t{{- range $value := .Values}}\n\t\t{{- with .Description}}\n\t\t\t{{.|prefixLines \"// \"}}\n\t\t{{- end}}\n\t\t{{$enum.GoType}}{{ .Name|toCamel }} {{$enum.GoType}} = {{.Name|quote}}\n\t{{- end }}\n\t)\n\n\tfunc (e {{.GoType}}) IsValid() bool {\n\t\tswitch e {\n\t\tcase {{ range $index, $element := .Values}}{{if $index}},{{end}}{{ $enum.GoType }}{{ $element.Name|toCamel }}{{end}}:\n\t\t\treturn true\n\t\t}\n\t\treturn false\n\t}\n\n\tfunc (e {{.GoType}}) String() string {\n\t\treturn string(e)\n\t}\n\n\tfunc (e *{{.GoType}}) UnmarshalGQL(v interface{}) error {\n\t\tstr, ok := v.(string)\n\t\tif !ok {\n\t\t\treturn fmt.Errorf(\"enums must be strings\")\n\t\t}\n\n\t\t*e = {{.GoType}}(str)\n\t\tif !e.IsValid() {\n\t\t\treturn fmt.Errorf(\"%s is not a valid {{.GQLType}}\", str)\n\t\t}\n\t\treturn nil\n\t}\n\n\tfunc (e {{.GoType}}) MarshalGQL(w io.Writer) {\n\t\tfmt.Fprint(w, strconv.Quote(e.String()))\n\t}\n\n{{- end }}\n",
	"object.gotpl":    "{{ $object := . }}\n\nvar {{ $object.GQLType|lcFirst}}Implementors = {{$object.Implementors}}\n\n// nolint: gocyclo, errcheck, gas, goconst\n{{- if .Stream }}\nfunc (ec *executionContext) _{{$object.GQLType}}(ctx context.Context, sel ast.SelectionSet) func() graphql.Marshaler {\n\tfields := graphql.CollectFields(ctx, sel, {{$object.GQLType|lcFirst}}Implementors)\n\tctx = graphql.WithResolverContext(ctx, &graphql.ResolverContext{\n\t\tObject: {{$object.GQLType|quote}},\n\t})\n\tif len(fields) != 1 {\n\t\tec.Errorf(ctx, \"must subscribe to exactly one stream\")\n\t\treturn nil\n\t}\n\n\tswitch fields[0].Name {\n\t{{- range $field := $object.Fields }}\n\tcase \"{{$field.GQLName}}\":\n\t\treturn ec._{{$object.GQLType}}_{{$field.GQLName}}(ctx, fields[0])\n\t{{- end }}\n\tdefault:\n\t\tpanic(\"unknown field \" + strconv.Quote(fields[0].Name))\n\t}\n}\n{{- else }}\nfunc (ec *executionContext) _{{$object.GQLType}}(ctx context.Context, sel ast.SelectionSet{{if not $object.Root}}, obj *{{$object.FullName}} {{end}}) graphql.Marshaler {\n\tfields := graphql.CollectFields(ctx, sel, {{$object.GQLType|lcFirst}}Implementors)\n\t{{if $object.Root}}\n\t\tctx = graphql.WithResolverContext(ctx, &graphql.ResolverContext{\n\t\t\tObject: {{$object.GQLType|quote}},\n\t\t})\n\t{{end}}\n\n\t{{if $object.IsConcurrent}} var wg sync.WaitGroup {{end}}\n\tout := graphql.NewOrderedMap(len(fields))\n\tinvalid := false\n\tfor i, field := range fields {\n\t\tout.Keys[i] = field.Alias\n\n\t\tswitch field.Name {\n\t\tcase \"__typename\":\n\t\t\tout.Values[i] = graphql.MarshalString({{$object.GQLType|quote}})\n\t\t{{- range $field := $object.Fields }}\n\t\tcase \"{{$field.GQLName}}\":\n\t\t\t{{- if $field.IsConcurrent }}\n\t\t\t\twg.Add(1)\n\t\t\t\tgo func(i int, field graphql.CollectedField) {\n\t\t\t{{- end }}\n\t\t\t\tout.Values[i] = ec._{{$object.GQLType}}_{{$field.GQLName}}(ctx, field{{if not $object.Root}}, obj{{end}})\n\t\t\t\t{{- if $field.ASTType.NonNull }}\n\t\t\t\t\tif out.Values[i] == graphql.Null {\n\t\t\t\t\t\tinvalid = true\n\t\t\t\t\t}\n\t\t\t\t{{- end }}\n\t\t\t{{- if $field.IsConcurrent }}\n\t\t\t\t\twg.Done()\n\t\t\t\t}(i, field)\n\t\t\t{{- end }}\n\t\t{{- end }}\n\t\tdefault:\n\t\t\tpanic(\"unknown field \" + strconv.Quote(field.Name))\n\t\t}\n\t}\n\t{{if $object.IsConcurrent}} wg.Wait() {{end}}\n\tif invalid { return graphql.Null }\n\treturn out\n}\n{{- end }}\n",
	"resolver.gotpl":  "//go:generate gorunpkg github.com/99designs/gqlgen\n\npackage {{ .PackageName }}\n\nimport (\n{{- range $import := .Imports }}\n\t{{- $import.Write }}\n{{ end }}\n)\n\ntype {{.ResolverType}} struct {}\n\n{{ range $object := .Objects -}}\n\t{{- if $object.HasResolvers -}}\n\t\tfunc (r *{{$.ResolverType}}) {{$object.GQLType}}() {{ $object.ResolverInterface.FullName }} {\n\t\t\treturn &{{lcFirst $object.GQLType}}Resolver{r}\n\t\t}\n\t{{ end -}}\n{{ end }}\n\n{{ range $object := .Objects -}}\n\t{{- if $object.HasResolvers -}}\n\t\ttype {{lcFirst $object.GQLType}}Resolver struct { *Resolver }\n\n\t\t{{ range $field := $object.Fields -}}\n\t\t\t{{- if $field.IsResolver -}}\n\t\t\tfunc (r *{{lcFirst $object.GQLType}}Resolver) {{ $field.ShortResolverDeclaration }} {\n\t\t\t\tpanic(\"not implemented\")\n\t\t\t}\n\t\t\t{{ end -}}\n\t\t{{ end -}}\n\t{{ end -}}\n{{ end }}\n",
	"server.gotpl":    "package main\n\nimport (\n{{- range $import := .Imports }}\n\t{{- $import.Write }}\n{{ end }}\n)\n\nconst defaultPort = \"8080\"\n\nfunc main() {\n\tport := os.Getenv(\"PORT\")\n\tif port == \"\" {\n\t\tport = defaultPort\n\t}\n\n\thttp.Handle(\"/\", handler.Playground(\"GraphQL playground\", \"/query\"))\n\thttp.Handle(\"/query\", handler.GraphQL({{.ExecPackageName}}.NewExecutableSchema({{.ExecPackageName}}.Config{Resolvers: &{{.ResolverPackageName}}.Resolver{}})))\n\n\tlog.Printf(\"connect to http://localhost:%s/ for GraphQL playground\", port)\n\tlog.Fatal(http.ListenAndServe(\":\" + port, nil))\n}\n",
}
