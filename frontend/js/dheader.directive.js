'use strick';

(function() {
    angular.module('dbus-profiler')
        .directive('dHeader', function() {
            return {
                restrict: 'EA',
                templateUrl: "templates/dheader.directive.html",
                controller: ['$scope', '$http', cc],
                compile: function() {
                    console.log("Caa")
                    return function() { console.log("Liiii")}
                }
            }

            function cc($scope, $http) {
                console.log("CC")
                $scope.update = function() {
                    $http({method: 'GET', url: '/dbus/api/info'}).success(function(data) {
                        $scope.info = data;
                    });
                }
                setInterval($scope.update, 1000);
                $scope.update()
            }
        });
})()
