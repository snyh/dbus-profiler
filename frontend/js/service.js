(function() {
    'use strict';
    angular.module('dbus-profiler')
        .factory('dapi', ['$http', function($http) {
            var dbus = "/dbus/api/"
            return {
                BuildHeaderInfo: function() {
                    return function() {
                        return get(dbus, "/info")
                    }
                },

                BuildIfcInfo: function(name) {
                    return function() {
                        return get(dbus, "interface?name="+name)
                    }
                },
                BuildMethodInfo: function(ifc, type, name) {
                    return function() {
                        return MethodInfos(ifc, type, name)
                    }
                },

                ConfigInfo: function() {
                    return $http.get("/config").then(function(resp) {
                        return resp.data.Enable
                    });
                },

                IfcInfos: ListIfcInfos,

                EnableAutoStart : function(v) {
                    var p = "f"
                    if (v) {
                        p = "t"
                    }
                    return $http.get("/config?enable=" + p).then(function(d) {
                    }, function(err) {
                        console.log(err)
                    })
                }
            }

            function get(base, name) {
                return $http.get(base + '/' + name).then(
                    function(resp) {
                        return resp.data
                    },
                    function(err) {
                        console.log("Errr on get ", url, err)
                    }
                );
            }
            function ListIfcInfos() {
                return get(dbus,"/main?1s")
                    .then(function(data) {
                        return data.map(function(d) {
                            return {
                                name: d.Ifc,
                                cost: d.TotalCost
                            }
                        })
                    })
            }
            function MethodInfos(ifcName, type, methodName) {
                return get(dbus, "interface?name="+ifcName).then(function(data) {
                    var v
                    switch (type) {
                    case "M":
                        v = data.Method[methodName]
                        break;
                    case "S":
                        v = data.Signal[methodName]
                        break;
                    case "P":
                        v = data.Property[methodName]
                        break;
                    }
                    return {
                        Ifc: ifcName,
                        Method: methodName,
                        Type: type,
                        Value: {
                            Total: v.Total,
                            Cost: v.Cost.map(function(d) { return d / 1000 / 1000; })
                        }
                    }
                })
            }
        }])
})()
