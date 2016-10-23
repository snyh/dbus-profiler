(function() {
    'use strick';
    angular.module('dbus-profiler')
        .directive('dHeader', function() {
            return {
                restrict: 'EA',
                templateUrl: "templates/dheader.directive.html",
                controller: ['$scope', 'dapi', cc],
            }

            function cc($scope, dapi) {
                $scope.update = function() {
                    dapi.get("/info", function(data){
                        $scope.info = data;
                    })
                }
                setInterval($scope.update, 1000);
                $scope.update()
            }
        });
})()
