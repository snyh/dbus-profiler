var duration = 0

function tick() {
    d3.json("/dbus/api/main?sort=name&since="+duration+"s", function(error, data) {
        if (error)
            return console.log(error);
        render(data, 1400, 760);
    });
    duration += 1
}
setInterval(tick, 1000);
tick();

var chart = d3.select('#chart')
    .append('svg')
    .attr('height', '100%')
    .attr('width', '100%')

const leftPadding = 48;
const bottomPadding = 20;
const iHeight = 40;


var tip = d3.tip()
    .attr('class', 'd3-tip')
    .offset([iHeight /2, 500])
    .html(function(d) {
        var c = format("total call ({}) , cost {}ms", d.TotalCall, (d.TotalCost/1000/1000.0));
        return format("<span style='color:red'>{}</span>", c)
    })

chart.call(tip);

function render(data, width, height) {
    
    var fill = d3.schemaCategory20c;
    min = d3.min(data, function(d) { return d.TotalCost; });
    max = d3.max(data, function(d) { return d.TotalCost; });

    var x = d3.scaleLinear()
        .range([0, height])
        .domain([min, max]);
    var nameScale = d3.scaleBand()
        .domain(data.map(function(o) { return o.Ifc })) 
        .range([0, height-bottomPadding],1,0.5)
    var tlScale = d3.scaleLinear()
        .domain([0, 60])
        .range([0, width])


    var yPosition = function(i, offset) {
        return i*iHeight + offset
    }

    chart.selectAll('.item').remove();
    chart.selectAll('.aix').remove();
    
    chart.append('g').call(d3.axisLeft(nameScale)).
        attr("class", "aix").
        attr('transform', format('translate({}, 0)', leftPadding))
    chart.append('g').call(d3.axisBottom(tlScale)).
        attr("class", "aix").
        attr('transform', format("translate(48, {})", height-bottomPadding));
    
    
    var group = chart.selectAll('.item').data(data);

    enter = group.enter().append('g').attr('class', 'item')        .on('mouseover', tip.show)
        .on('mouseout', tip.hide)


    enter.append('rect')
        .attr('height', iHeight/2)
        .attr('transform', function(d, i) { return format("translate({},{})", leftPadding+1, yPosition(i,0)); })
        .attr('width', function(d) {
            return x(d.TotalCost);
        })
        .style('fill', 'rgba(0,200,10,0.6)').style('stroke', 'black')

    enter.append('text')
        .attr('transform', function(d, i) { return format("translate(100,{})", yPosition(i,15)); })
        .text(function(d) {
            return d.Ifc
        });

    group.exit().remove();
}

