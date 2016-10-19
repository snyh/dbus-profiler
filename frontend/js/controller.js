
var app = angular.module('dbus-profiler', []);
app.controller('HeaderCtrl', function($scope, $http) {
    setInterval(function() {
        $http({method: 'GET', url: '/dbus/api/info'}).success(function(data) {
            $scope.info = data;
        });
    }, 1000);
});


app.controller('DetailCtrl', function($scope,$http) {
});
