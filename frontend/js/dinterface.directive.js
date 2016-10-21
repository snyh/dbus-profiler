'use strict';

(function(){
    angular.module('dbus-profiler')
        .directive('ifcName', function() {
            return {
                controller: function($scope) {}
            }
        })
        .directive('dInterface', function() {
            return {
                replace: true,
                require: '^ifcName',
                scope: {
                    ifcName: '@'
                },
                restrict: 'EA',
                templateUrl: "templates/dinterface.directive.html",
                controller: ['$scope', '$http', ic],
                link: link,
            }
        });

    function link(scope, iElement, iAttrs) {
        console.log("LINK>>>>>>>>>>>>")
        setInterval(scope.update, 1000)
        scope.update()

        scope.$watch('detailName', function(newVal) {
            if (scope.info && newVal) {
                var v = scope.info.Method[newVal]
                var root = iElement[0].querySelector(".dinterface_chart_container")
                root.innerHTML = ""
                draw_detail(root, v.Cost.map(function(d) {return d / 1000 /1000;}), iAttrs)
            }
        })
    }
    
    function ic($scope, $http) {
        $scope.switchM = function(n) {
            $scope.detailName = n
        }
        $scope.update = function() {
            $scope.getData($scope.ifcName || "org.freedesktop.DBus")
        }
        
        $scope.getData = function(name) {
            $http({method: 'GET', url: '/dbus/api/interface?name='+name}).success(function(data) {
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
    }
    
   function draw_detail(root, data, opts)
    {
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
