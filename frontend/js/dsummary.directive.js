(function() {
    'use strict';
    angular.module('dbus-profiler')
        .component('dSummary', {
            controller: ['$scope', '$element', '$attrs', cc],
            template: '<ui-view></ui-view> <div id=chart></div>'
        });

    function cc(scope, iElement, iAttrs) {
        scope.$watch("data", function(newVal) {
            if (newVal) {
                var root = iElement[0]
                root.innerHTML = ""
                render(scope, root, newVal, iAttrs)
            }
        })
        window.onresize = scope.tick;

        scope.tick = function() {
            d3.json(format("/dbus/api/main?top={}&since={}s", MaxServer, MaxSecond), function(error, data) {
                if (error)
                    return console.log(error);
                scope.data = data.sort(function(a,b) {
                    return a.Ifc > b.Ifc
                });
            });
        }
        setInterval(scope.tick, 1000);
        scope.tick()
    }

    const leftPadding = 48;
    const bottomPadding = 48;
    const MaxSecond = 60;
    const MaxServer = 10;

    function render(scope, root, data, opts) {
        var width = opts.width || 800,
            height = opts.height || 600,
            numberServer = data.length,
            iHeight = height / numberServer - 10

        var yPosition = d3.scaleLinear()
            .domain([0, numberServer])
            .range([0, height-bottomPadding],1,0.5)

        var x = d3.scaleLinear()
            .domain(d3.extent(data.map(function(d){return d.TotalCost})))
            .range([10, height]);

        var svg = d3.select(root)
            .append('svg:svg')
            .attr('width', width)
            .attr('height', height)
            .append('g')

        svg.selectAll("*").remove();

        var update = svg.selectAll('.item').data(data);

        update.exit().remove();

        var enter = update.enter()
            .append('a').attr('href', function(d) { return "/static/#/ifcs/" + d.Ifc })
            .append('g')
            .classed('item', true)

        enter.append('rect')
            .attr('height', iHeight/2)
            .attr('transform', function(d, i) { return format("translate({},{})", leftPadding, yPosition(i)+iHeight/3.0) })
            .attr('width', function(d) {
                return x(d.TotalCost);
            })
            .attr('fill', function(d, i) { return d3.schemeCategory10[i] })
            .attr('stroke', "black")

        enter.append('text')
            .attr('transform', function(d, i) { return format("translate(100,{})", yPosition(i) + iHeight/3.0 + 16) })
            .text(function(d) {
                return format("{} total call ({}) , cost {}ms", d.Ifc, d.TotalCall, (d.TotalCost/1000/1000.0));
            });


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
