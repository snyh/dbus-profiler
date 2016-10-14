var values = [0,1,2,3,4,5,6,7,8,9,2,3,4,5,5,5,5,5,1,1,7,7,7,6,4,8];

var width = 960,
    height = 500,
    padding=6;

var x = d3.scale.linear()
    .domain([0, 10])
    .range([0, width/2]);
var data = d3.layout.histogram()
  	.bins(5)
(values);

var xAxis = d3.svg.axis()
    .scale(x)
    .orient("bottom");

var y = d3.scale.linear()
    .domain([0, 10])
    .range([height, 0]) ;
var yAxis = d3.svg.axis()
    .scale(y)
    .orient("left");
var svg = d3.select("body").append("svg")
    .attr("width", width+30)
    .attr("height", height+30)
    .append("g")
    .attr("transform", "translate(" + padding*4 + "," + padding + ")");

var bar = svg.selectAll(".bar")
    .data(data)
    .enter().append("g")
    .attr("class", "bar")
    .attr("transform", function(d,i) { return "translate(" + x(i*2) + "," + y(d.y) + ")"; });
bar.append("rect")
 	.attr("x",padding)
    .attr("width", x(data[0].dx))
    .attr("height", function(d) { return height - y(d.y); });
bar.append("text")
    .attr("dy", ".75em")
    .attr("y", 6)
    .attr("x", x(data[0].dx) / 2)
    .attr("text-anchor", "middle")
    .text(function(d) { return d.y; })
    .attr("fill", "#fff");
svg.append("g")
    .attr("class", "x axis")
    .attr("transform", "translate("+padding+"," + height + ")")
    .call(xAxis);
svg.append("g")
    .attr("class", "y axis")
    .attr("transform", "translate("+padding+"," + 0 + ")")
    .call(yAxis);
