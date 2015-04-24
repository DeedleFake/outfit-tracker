package ot

import (
	"text/template"
)

const (
	pageMain = `<html>
	<head>
		<link rel='stylesheet' href='/assets/style.css' />
	</head>
	<body>
		<form method='get' action='/view'>
			Server:
			<select name='server'>
				<!--<option value='19'>Jaeger</option>-->
				<option value='25'>Briggs</option>
				<option value='13'>Cobalt</option>
				<option value='10'>Miller</option>
				<option value='1'>Connery</option>
				<option value='17'>Emerald</option>
			</select>

			Faction:
			<select name='faction'>
				<!--<option value='0'>Nanite Systems</option>-->
				<option value='1'>Vanu Sovereignty</option>
				<option value='2'>New Conglomerate</option>
				<option value='3'>Terran Republic</option>
			</select>

			<input type='submit' value='View' />
		</form>
</html>`

	pageView = `<html>
	<head>
		<link rel='stylesheet' href='/assets/style.css' />

		<script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>
		<script src="http://cytoscape.github.io/cytoscape.js/api/cytoscape.js-latest/cytoscape.min.js"></script>
		<script type='application/javascript' src='/assets/view.js'></script>
	</head>
	<body>
		<div id='header'><h2>{{.server}}: {{.faction}}</h2></div>

		<div id='graph'></div>

		<div style='display:none;'>
			<div id='server'>{{.serverID}}</div>
			<div id='faction'>{{.factionID}}</div>
		</div>
	</body>
</html>`
)

var (
	tmpl = template.New("main")
)

func init() {
	template.Must(tmpl.Parse(pageMain))
	template.Must(tmpl.New("view").Parse(pageView))
}
