package docs

import _ "embed"

//go:embed swagger.yaml
var swaggerYAML []byte

//go:embed swagger.json
var swaggerJSON []byte

//go:embed swagger-ui.css
var swaggerUICSS []byte

//go:embed swagger-ui-bundle.js
var swaggerUIJS []byte

//go:embed swagger-ui.html
var swaggerUIHTML []byte
