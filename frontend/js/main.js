var duration = 0

function tick() {
    d3.json("/dbus/api/main?sort=cost&since="+duration+"s", function(error, json) {
        if (error)
            return console.log(json);
        render(json);
    });
    duration += 1
}
setInterval(tick, 1000);
tick();


//setInterval(tick, 5000);
var chart = d3.select('#chart')
    .append('svg')
    .attr('width', 1920)
    .attr('height', 1080);

function render(data) {
    var fill = d3.schemaCategory20c;
    min = d3.min(data, function(d) { return d.Cost; });
    max = d3.max(data, function(d) { return d.Cost; });

    var x = d3.scaleLinear()
        .range([0, 1000])
        .domain([min, max]);

    var iHeight = 40;
    var yPosition = function(i, offset) {
        return i*iHeight + offset
    }

    chart.selectAll('g').remove();
    
    var group = chart.selectAll('g').data(data);
    
    enter = group.enter().append('g');

    enter.append('rect')
        .attr('height', iHeight/2)
        .attr('class', 'item')
        .attr('transform', function(d, i) { return "translate(0,"+yPosition(i,0)+")"; })
        .attr('width', function(d) {
            return x(d.Cost);
        })
        .style('fill', 'blue').style('stroke', 'black')


    enter.append('text')
        .attr('class', 'item')
        .attr('transform', function(d, i) { return "translate(100,"+yPosition(i,15)+")"; })
        .text(function(d) {
            return "" + d.Ifc + "(" + d.RCs.length + ")" + " " + (d.Cost /1000 /1000.0) + "ms";
        });

    group.exit().remove();
}

