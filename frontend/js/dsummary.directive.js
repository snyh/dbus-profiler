(function() {
    'use strict';
    angular.module('dbus-profiler')
        .directive('dSummary', function() {
            return {
                restrict: 'EA',
                link: link,
                controller: ['$scope', 'dapi', cc],
            }
        });

    function link(scope, iElement, iAttrs) {
        scope.$watch("data", function(newVal) {
            if (newVal) {
                var root = iElement[0]
                root.innerHTML = ""
                render(scope, root, newVal, iAttrs)
            }
        })
        window.onresize = scope.tick;
        setInterval(scope.tick, 1000);
        scope.tick()
    }

    function cc(scope) {
        scope.tick = function() {
            d3.json(format("/dbus/api/main?top={}&since={}s", MaxServer, MaxSecond), function(error, data) {
                if (error)
                    return console.log(error);
                scope.data = data;
            });
        }
    }

    var yPosition;
    var iHeight = 40;
    const leftPadding = 48;
    const bottomPadding = 48;
    const MaxSecond = 60;
    const MaxServer = 10;
    var numberServer = 0;

    function render(scope, root, data, opts) {
        var width = opts.width || 800,
            height = opts.height || 600;

        if (numberServer != data.length) {
            numberServer = data.length
            iHeight = height / (numberServer) - 10

            yPosition = d3.scaleLinear()
                .domain([0, numberServer])
                .range([0, height-bottomPadding],1,0.5)
        }

        var x = d3.scaleLinear()
            .domain(d3.extent(data.map(function(d){return d.TotalCost})))
            .range([10, height]);

        var svg = d3.select(root)
            .append('svg:svg')
            .attr('width', width)
            .attr('height', height)
            .append('g')

        svg.selectAll("*").remove();

        var group = svg.selectAll('.item').data(data);

        var enter = group.enter().append('g')
            .classed('item', true)
            .on("click", function(d, i, ele) {
                    scope.switchIfc(d.Ifc)
            })

        enter.append('rect')
            .attr('height', iHeight/2)
            .attr('transform', function(d, i) { return format("translate({},{})", leftPadding, yPosition(i)+iHeight/3.0) })
            .attr('width', function(d) {
                return x(d.TotalCost);
            })
            .style('fill', 'rgba(0,200,10,0.6)').style('stroke', 'black')

        enter.append('text')
            .attr('transform', function(d, i) { return format("translate(100,{})", yPosition(i) -25) })
            .attr('fill', function(d, i) { return d3.schemeCategory10[i] })
            .text(function(d) {
                return format("{} total call ({}) , cost {}ms", d.Ifc, d.TotalCall, (d.TotalCost/1000/1000.0));
            });

        group.exit().remove();

        renderAix(svg, data, width, height)
        renderPath(scope.ifcName, svg, data, width, height)
    }

    function renderAix(svg, data, width, height) {
        var yScale = d3.scaleLinear()
            .domain([50, 0])
            .range([0, height-bottomPadding])

        var tlScale = d3.scaleLinear()
            .domain([0, MaxSecond])
            .range([0, width])

        svg.selectAll('.aix').remove();

        svg.append('g').call(d3.axisLeft(yScale))
            .classed('aix', true)
            .attr('transform', format('translate({}, 0)', leftPadding))

        svg.append('g').call(d3.axisBottom(tlScale))
            .classed('aix', true)
            .attr('transform', format("translate(48, {})", height-bottomPadding));
    }


    function renderPath(ifcName, svg, rawdata, width, height) {
        var tlScale = d3.scaleLinear()
            .domain([0, MaxSecond])
            .range([bottomPadding ,width-leftPadding])

        var data = rawdata.map(function(d){return d.CostDetail;})

        var y = d3.scaleLinear()
            .domain([0, 1000*1000*1000])
            .range([bottomPadding, height])

        var fn = d3.line().curve(d3.curveBasis)
            .x(function(d, j) { return width - tlScale(j) })
            .y(function(d, j) {
                return height - y(d)
            })

        svg.selectAll("path").remove()
        svg.selectAll(".item path").data(data)
            .enter().append('path')
            .attr('d', fn)
            .attr('stroke', function(d, i) { return d3.schemeCategory10[i] })
            .attr('stroke-width', 3)
            .attr('fill', 'none')
    }


})()
