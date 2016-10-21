function draw_detail(root, data, width, height)
{
    var graphics = d3.select(root).append("g")
    graphics.select("g.parent").remove();

    var leftPadding = 25
    var bottomPadding = 25
    
    var max_height = height - leftPadding;
    var max_width = width - bottomPadding;

    var total = data.length
    var hist = d3.histogram()
    


    var y = d3.scaleLinear()
        .range([0, max_height])

    var yTick = d3.scaleLinear()
        .range([0, data.length])
        .domain(0, max_height)
    
    var x = d3.scaleLinear()
        .domain(d3.extent(data))
        .range([0,max_width])

    graphics.append('g').call(d3.axisBottom(x).tickFormat(d3.format(",.1f")))
        .attr('transform', format('translate({}, {})', leftPadding, max_height));

    graphics.append('g').call(d3.axisLeft(yTick))
        .attr('transform', format('translate({}, 0)', leftPadding))

    update = graphics.append('g')
        .attr('transform', format('translate({}, 0)', leftPadding))
        .selectAll("rect").data(hist(data))

	update.enter()
		.append("rect")
        .attr('width', function(d) { return x(d.x1) - x(d.x0); })
        .attr('height', function(d) { return max_height - y(d.length/total)})
        .attr('transform', function(d) { return format("translate({}, {})", x(d.x0), y(d.length/total)) })
        .attr('stroke', 'blue')
		.attr("fill","steelblue");

    update.exit().remove()

    return;    
}




function renderAix(data, width, height) {
    var nameScale = d3.scaleBand()
        .domain(data.map(function(d) { return d.Ifc })) 
        .range([0, height-bottomPadding],1,0.5)
    
    var tlScale = d3.scaleLinear()
        .domain([0, MaxSecond])
        .range([0, width])
    
    chart.selectAll('.aix').remove();
    
    chart.append('g').call(d3.axisLeft(nameScale))
        .classed('aix', true)
        .attr('transform', format('translate({}, 0)', leftPadding))
    
    chart.append('g').call(d3.axisBottom(tlScale))
        .classed('aix', true)
        .attr('transform', format("translate(48, {})", height-bottomPadding));
}


function renderPath(data, width, height) {
    var min = d3.min(data, function(d) { return d.TotalCost; });
    var max = d3.max(data, function(d) { return d.TotalCost; });


    var tlScale = d3.scaleLinear()
        .domain([0, MaxSecond])
        .range([bottomPadding ,width-leftPadding])

    chart.selectAll("path").remove()
    chart.selectAll(".item path").data(data.map(function(d){return d.CostDetail;}))
        .enter().append('path')
        .attr('d', function(d, i) {
            var fn = d3.line().curve(d3.curveBasis)
                .x(function(d, j) { return width - tlScale(j) })
                .y(function(d, j) {
                    var y = d3.scaleLinear()
                        .domain([0, 1000 * 1000* 10])
                        .range([bottomPadding, height])
                    return height - y(d)
                })
            return fn(d)
        })
        .attr('stroke', function(d, i) { return d3.schemeCategory10[i] })
        .attr('stroke-width', 3)
        .attr('fill', 'none')
 }

function render(data, width, height) {
    min = d3.min(data, function(d) { return d.TotalCost; });
    max = d3.max(data, function(d) { return d.TotalCost; });

    var x = d3.scaleLinear()
        .domain([min, max])
        .range([10, height]);
    
    chart.selectAll('.item').remove();    

    var group = chart.selectAll('.item').data(data);

    enter = group.enter().append('g')
        .classed('item', true)
        .on("click", function(d) { SwitchInterface(d.Ifc) })

    enter.append('rect')
        .attr('height', iHeight/2)
        .attr('transform', function(d, i) { return format("translate({},{})", leftPadding, yPosition(i)+iHeight/3.0) })
        .attr('width', function(d) {
            return x(d.TotalCost);
        })
        .style('fill', 'rgba(0,200,10,0.6)').style('stroke', 'black')

    enter.append('text')
        .attr('transform', function(d, i) { return format("translate(100,{})", yPosition(i) -25) })
        .text(function(d) {
            return format("{} total call ({}) , cost {}ms", d.Ifc, d.TotalCall, (d.TotalCost/1000/1000.0));
        });

    group.exit().remove();

    renderAix(data, width, height)
    renderPath(data, width, height)
}
