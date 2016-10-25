(function() {
    'use strick';
    angular.module('dbus-profiler')
        .component('dHeader', {
            templateUrl: "templates/dheader.directive.html",
            controller: ['$scope', 'dapi', function($scope, dapi) {
                var fetchFn = dapi.BuildHeaderInfo()


                $scope.enableAutostart = function() {
                    dapi.EnableAutoStart(!!$scope.autoStart)
                }

                var update = function() {

                    fetchFn().then(function(d) {
                        $scope.info = d
                    })

                    dapi.ConfigInfo().then(function(v) {
                        $scope.autoStart = v
                    })

                }
                setInterval(update, 1000);
                update();
            }]
        })
})()
