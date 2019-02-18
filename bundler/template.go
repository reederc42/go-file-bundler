package bundler

const bundleTemplate = `//generated file
package {{.Package}}

var {{.Name}} = {{.Value}}
`
