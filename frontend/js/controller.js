var app = angular.module('dbus-profiler', ['smart-table']);
//var app = angular.module('dbus-profiler', []);

app.filter('sum', function() {
    return function(d) {
        var sum = 0
        angular.forEach(d, function(i) { sum = sum + i})
        return sum
    }
});

app.filter('percentage', ['$filter', function($filter) {
    return function(input, all, decimals) {
        return $filter('number')(input * 100.0 / all, decimals) + '%';
    }
}]);

app.controller('HeaderCtrl', ['$scope', '$http', function($scope, $http) {
    setInterval(function() {
        $http({method: 'GET', url: '/dbus/api/info'}).success(function(data) {
            $scope.info = data;
        });
    }, 1000);
}]);

app.controller('DetailCtrl', ['$scope', '$http', function($scope, $http) {
    var rows = []
    $scope.rows = [
        {name:"ABBC", type:"M", call:123, cost:1111},
        {name:"BBCC", type:"M", call:3932, cost:11},
        {name:"ZZZZ", type:"S", call:123, cost:10},
    ]
    $scope.update = function() {
        if ($scope.IfcName != "") {
            $http({method: 'GET', url: '/dbus/api/interface?name='+$scope.IfcName}).success(function(data) {
                rows = []
                $scope.info = data
                angular.forEach($scope.info.Method, function(v, name) {
                    rows.push({name:name, type:"M", call: v.Total, cost: v.Cost.reduce(function(a, i) { return a + i})});
                })
                angular.forEach($scope.info.Property, function(v, name) {
                    rows.push({name:name, type:"P", call: v.Total, cost: v.Cost.reduce(function(a, i) { return a + i})});
                })
                angular.forEach($scope.info.Signal, function(v, name) {
                    rows.push({name:name, type:"S", call: v.Total, cost: 0});
                })
                $scope.rows = rows;
            });

        } else {
            console.log("hehe...")
            $scope.info = null
        }
    }

    $scope.SwitchInterface = function(name) {
        $scope.IfcName = name
        $scope.update()
    }
    


    $scope.switchMethod = function(name) {
        $scope.detailName = name
        v = $scope.info.Method[name]
        root = document.querySelector("#method_detail_chart svg")
        root.innerHTML = "";
        draw_detail(root, v.Cost.map(function(d) {return d / 1000 /1000;}), 250, 200)
    }
    $scope.switchProperty = function(name) {
        $scope.detailName = name

        v = $scope.info.Method[name]
        root = document.querySelector("#method_detail_chart svg")
        root.innerHTML = "";
        draw_detail(root, v.Cost.map(function(d) {return d / 1000 /1000;}), 250, 200)
    }
    $scope.switchSignal = function(name) {
        $scope.detailName = name
        root = document.querySelector("#method_detail_chart svg")
        root.innerHTML = "";
        console.log("SwitchSignalTo" + name)
    }

    $scope.SwitchInterface("org.freedesktop.DBus")
    setInterval($scope.update, 1000)
}]);


app.controller('SummaryCtrl', ['$scope', function($scope) {
    $scope.tick = function() {
        d3.json(format("/dbus/api/main?top={}&since={}s", MaxServer, MaxSecond), function(error, data) {
            if (error)
                return console.log(error);
            
            var width = header.clientWidth;
            var height =  document.body.clientHeight - header.clientHeight;
            
            if (numberServer != data.length) {
                numberServer = data.length
                iHeight = height / (numberServer) - 10
                
                yPosition = d3.scaleLinear()
                    .domain([0, numberServer])
                    .range([0, height-bottomPadding],1,0.5)
            }
            render(data, width, height);
        });
    }

    setInterval($scope.tick, 1000);
    window.onresize = $scope.tick;
    $scope.tick()
}]);
