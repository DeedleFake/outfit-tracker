$(document).ready(function() {
	var server = $('#server').html();
	var faction = $('#faction').html();

	$.getJSON('/data?server=' + server + '&faction=' + faction, function(data) {
		cytoscape({
			container: document.getElementById('graph'),

			style: cytoscape.stylesheet()
				.selector('node').css({
					'content': 'data(id)'
				})
				.selector('edge').css({
					'target-arrow-shape': 'triangle',
					'width': 3,
					'line-color': '#DDDDDD',
					'target-arrow-color': '#999999'
				}),

			elements: data,

			layout: {
				name: 'circle',
				directed: true
			}
		});
	});
});
