$(document).ready(function() {
	var server = $('#server').html();
	var faction = $('#faction').html();

	$.getJSON('/data?server=' + server + '&faction=' + faction, function(data) {
		cytoscape({
			container: document.getElementById('graph'),

			style: cytoscape.stylesheet()
				.selector('node').css({
					'content': 'data(label)',
					'text-valign': 'center',
					'color': '#FFFFFF',
					'text-outline-width': 3,
					'text-outline-color': '#888888',
					'width': 'data(size)',
					'height': 'data(size)'
				})
				.selector('edge').css({
					'target-arrow-shape': 'triangle',
					'width': 3,
					'line-color': '#DDDDDD',
					'target-arrow-color': '#999999'
				}),

			elements: data,

			layout: {
				name: 'concentric',
				concentric: function() {
					return this.data('size') * 100;
				},
				levelWidth: function(nodes) {
					return 100;
				},
				avoidOverlap: true,
				directed: true,
				padding: 10
			}
		});
	});
});
