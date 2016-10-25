(function(){
    'use strict';
    angular.module('dbus-profiler')
        .component('dIfcList', {
            bindings: {
                infos: '<'
            },
            templateUrl: "templates/difclist.html"
        })
        .component('dInterface', {
            bindings: {
                fetchFn: '<'
            },
            templateUrl: "templates/dinterface.directive.html",
            controller: ['$scope', '$element', '$attrs', ic]
        })
        .component('dMethod', {
            bindings: {
                fetchFn: '<'
            },
            templateUrl: "templates/dmethod.html",
            controller: ['$scope', '$element', '$attrs', function($scope, $element, $attrs) {
                var fetchFn = $scope.$ctrl.fetchFn
                var chart = $element[0].querySelector(".chart")

                $scope.$watch('cost', function(newVal) {
                    if (newVal) {
                        $scope.max = d3.max(newVal)
                        $scope.min = d3.min(newVal)
                        $scope.mean = d3.mean(newVal)
                        $scope.median = d3.median(newVal)
                        $scope.deviation = d3.deviation(newVal)

                        if ($scope.type == "S") {
                            chart.innerHTML = ""
                        } else {
                            draw_detail(chart, newVal, $attrs)
                        }
                    }
                })

                var update = function() {
                    fetchFn().then(function(d) {
                        $scope.ifc = d.Ifc
                        $scope.name = d.Method
                        $scope.type = d.Type
                        $scope.total = d.Value.Total
                        $scope.cost = d.Value.Cost
                    })
                }
                update()
                setInterval(update, 1000)
            }]
        })

    function ic($scope, $element, $attrs) {
        var fetchFn = $scope.$ctrl.fetchFn
        $scope.update = function() {
            fetchFn().then(function(data) {
                var rows = []
                angular.forEach(data.Method, function(v, name) {
                    rows.push({name:name, type:"M", call: v.Total, cost: v.Cost.reduce(function(a, i) { return a + i})});
                })
                angular.forEach(data.Property, function(v, name) {
                    rows.push({name:name, type:"P", call: v.Total, cost: v.Cost.reduce(function(a, i) { return a + i})});
                })
                angular.forEach(data.Signal, function(v, name) {
                    rows.push({name:name, type:"S", call: v.Total, cost: 0});
                })

                $scope.rows = rows;
                $scope.info = data
            });
        }
        setInterval($scope.update, 1000)
        $scope.update()
    }

    function draw_detail(root, data, opts)
    {
        root.innerHTML = ""
        var width = opts.width || 200,
            height = opts.height || 250,
            leftPadding = opts.leftPadding || 30,
            bottomPadding = opts.bottomPadding || 25;

        var svg = d3.select(root)
            .append("svg")
            .attr('width', width)
            .attr('height', height)
            .append('g')
            .attr('transform', format("translate({}, {})", leftPadding, bottomPadding));

        svg.selectAll("*").remove();


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

        svg.append('g').call(d3.axisBottom(x).tickFormat(d3.format(",.1f")))
            .attr('transform', format('translate({}, {})', leftPadding, max_height));

        svg.append('g').call(d3.axisLeft(yTick))
            .attr('transform', format('translate({}, 0)', leftPadding))

        var update = svg.append('g')
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
})()
