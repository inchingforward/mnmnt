var API_KEY = 'pk.eyJ1IjoibWphbmdlciIsImEiOiJjazN6NHZlNHkwMjZiM2tudzRpN3FyNzc0In0.avUDs9ardvviib8L8HsMSA';
var tileLayer = L.tileLayer('https://api.mapbox.com/styles/v1/{id}/tiles/{z}/{x}/{y}?access_token={accessToken}', {
    attribution: 'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, <a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
    maxZoom: 18,
    //id: 'mapbox/streets-v11',
    id: 'mapbox/light-v10',
    tileSize: 512,
    zoomOffset: -1,
    accessToken: API_KEY
});

var RecentMemories = RecentMemories || (function() {
    var infoWindow, map, bounds, markers, prevMarker, group, ui;    
        
    function init() {
        map = L.map('map').setView([38.677811, -90.419197], 13);

        tileLayer.addTo(map);
        
        markers = {};
    }

    function showMarkerBubble(marker) {
        var bubbles = ui.getBubbles();
        for (var i=0; i < bubbles.length; i++) {
            ui.removeBubble(bubbles[i]);
        }

        var bubble =  new H.ui.InfoBubble(marker.getGeometry(), {
            content: marker.getData()
        });

        ui.addBubble(bubble);
    }

    function addMemory(latitude, longitude, memoryId, memoryTitle) {
        var marker = L.marker([latitude, longitude]);
        marker.bindPopup('<a href="/memories/' + memoryId + '">' + memoryTitle + '</a>');
        marker.addTo(map);

        markers[memoryId] = marker;

        if (!bounds) {
            bounds = new L.LatLngBounds(); 
        }

        bounds.extend(marker.getLatLng());

        map.fitBounds(bounds);
    }

    function showInfoWindowForMemory(memoryId) {
        var marker = markers[memoryId];
        marker.openPopup();
    }

    return {
        showInfoWindowForMemory: showInfoWindowForMemory,
        addMemory: addMemory, 
        init: init
    }
})();

var AddMemory = AddMemory || (function() {
    var map, marker, geocoder;

    var delay = function() {
        var timeout = 0;
        return function(callback, ms) {
            clearTimeout(timeout);
            timeout = setTimeout(callback, ms);
        };
    }();

    function init() {
        map = new mapboxgl.Map({
            container: 'address_search_map',
            style: 'mapbox://styles/mapbox/light-v10',
            center: [-90.1994, 38.6270],
            zoom: 13
        });

        geocoder = new MapboxGeocoder({
            accessToken: mapboxgl.accessToken,
            mapboxgl: mapboxgl,
            placeholder: "1. Search for a place or address" 
        });

        //document.getElementById('geocoder').appendChild(geocoder.onAdd(map));
        
        $("#address_text").keyup(function() {
            delay(findAddress, 1500);
        });

        $("#address_text").keypress(function (e) {
            if ((e.which && e.which == 13) || (e.keyCode && e.keyCode == 13)) {
                findAddress();
            }
        });

        // geocoder = new MapboxGeocoder({ accessToken: mapboxgl.accessToken });
        geocoder.on('result', function(result) {
            console.log(result);
            console.log(result.result.center);
        });

        console.log(geocoder);
        console.log(geocoder.query);
    }

    function findAddress() {
        var addressText = $("#address_text").val();
        if (!addressText) {
            alert("Please enter an address to search for.");
            return;
        }

        if (marker) {
            marker.setMap(null);
        }

        console.log("gonna query...")
        //geocoder.query(addressText);
        //geocoder.inputString = addressText;
        geocoder._geocode(addressText);
        console.log("queried");

        // var latitude = $("#latitude");
        // var longitude = $("#longitude");
        // var submitButton = $("#submit_button");
        // var geocoder = new google.maps.Geocoder();

        // $("#map_spinner").show();


        // geocoder.geocode({address: addressText}, function(results, status) {
        //     $("#map_spinner").hide();

        //     if (status == google.maps.GeocoderStatus.OK) {
        //         map.setCenter(results[0].geometry.location);

        //         marker = new google.maps.Marker({
        //             map: map,
        //             position: results[0].geometry.location
        //         });

        //         latitude.val(results[0].geometry.location.lat());
        //         longitude.val(results[0].geometry.location.lng());

        //         submitButton.attr("disabled", false);
        //     } else {
        //         latitude.val("");
        //         longitude.val("");

        //         submitButton.attr("disabled", "disabled");
        //     }
        // });
    }

    function verifyForm() {
        var latitude = $("#latitude").val();
        var longitude = $("#longitude").val();

        if (!latitude || !longitude) {
            alert("Please search and find an address.");
            $("#address_text").focus();
            return false;
        }

        var titleField = $("#title");
        if (!titleField.val()) {
            alert("Please enter a title.");
            titleField.focus();
            return false;
        }

        var detailsField = $("#details");
        if (!detailsField.val()) {
            alert("Please enter some details.");
            detailsField.focus();
            return false;
        }

        return true;
    }

    return {
        init: init,
        findAddress: findAddress
    }
})();

var EditMemory = AddMemory || (function() {
    var map, marker;

    function init(latitude, longitude) {
        var lngLat = [lng, lat];

        var map = new mapboxgl.Map({
            container: 'edit_memory_map',
            style: 'mapbox://styles/mapbox/light-v10',
            zoom: 13,
            center: lngLat
        });
        
        var marker = new mapboxgl.Marker()
            .setLngLat(lngLat)
            .addTo(map);
    }

    function verifyForm() {
        var detailsField = $("#details");
        if (!detailsField.val()) {
            alert("Please enter some details.");
            detailsField.focus();
            return false;
        }

        return true;
    }

    return {
        init: init,
    }
})();

var MemoryDetails = MemoryDetails || (function() {
    function showMap(lat, lng, title) {
        map = L.map('memory_map').setView([lat, lng], 13);

        tileLayer.addTo(map);
        
        L.marker([lat, lng]).addTo(map);
    }

    return {
        showMap: showMap
    }
})();