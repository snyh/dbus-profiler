

function tick() {
    d3.json("/dbus/api/main", function(error, json) {
        if (error)
            return console.log(json);
        draw(json);
    });
}

//setInterval(tick, 5000);
var chart = d3.select('#chart')
    .append('svg')
    .attr('width', 1920)
    .attr('height', 1080);

function draw(data) {
    var fill = d3.scale.category20();
    min = d3.min(data, function(d) { return d.Cost; });
    max = d3.max(data, function(d) { return d.Cost; });

    var x = d3.scale.linear()
        .range([0, 1000])
        .domain([min, max]);

    var iHeight = 40;

    chart.selectAll('g').remove();
    
    var group = chart.selectAll('g').data(data);
    
    enter = group.enter().append('g');

    enter.append('rect')
        .attr('class', 'item')
        .attr('transform', function(d, i) { return "translate(0," + i * iHeight +")"; })
        .style({'fill':'blue', 'stroke':'black'})
        .attr('width', function(d) {
            return x(d.Cost);
        })
        .attr('height', iHeight)


    enter.append('text')
        .attr('class', 'item')
        .attr('transform', function(d,i) { return "translate(0," + i * iHeight +")"; })
        .text(function(d) {
            return "" + d.Ifc + "(" + d.RCs.length + ")";
        });

    group.exit().remove();
}

setInterval(tick, 1000);
tick();
