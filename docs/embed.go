package docs

import _ "embed"

//go:embed swagger.yaml
var Swagger_yaml []byte

//go:embed swagger.json
var Swagger_json []byte

//go:embed swagger-ui.css
var Swagger_ui_css []byte

//go:embed swagger-ui-bundle.js
var Swagger_ui_js []byte

//go:embed swagger-ui.html
var Swagger_ui_html []byte
