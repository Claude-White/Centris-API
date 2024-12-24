import fs from "fs";

const delay = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

(async () => {
    fs.readFile("test.json", "utf-8", async (err, data) => {
        if (err) {
            console.error("Error reading the file:", err);
            return;
        }

        const jsonData = JSON.parse(data);

        for (const item of jsonData) {
            // Fetch the markers
            try {
                const res = await fetch(
                    "https://www.centris.ca/api/property/map/GetMarkers",
                    {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                        },
                        body: JSON.stringify({
                            zoomLevel: 11,
                            mapBounds: {
                                SouthWest: {
                                    Lat: item.mapBounds.SouthWest.Lat,
                                    Lng: item.mapBounds.SouthWest.Lng,
                                },
                                NorthEast: {
                                    Lat: item.mapBounds.NorthEast.Lat,
                                    Lng: item.mapBounds.NorthEast.Lng,
                                },
                            },
                        }),
                    }
                );

                if (!res.ok) {
                    console.error(
                        "Error fetching markers:",
                        res.status,
                        res.statusText
                    );
                    continue;
                }

                const data = await res.json();
                console.log("[");

                for (const marker of data.d.Result.Markers) {
                    // Fetch marker info
                    try {
                        const res = await fetch(
                            "https://www.centris.ca/property/GetMarkerInfo",
                            {
                                method: "POST",
                                headers: {
                                    "Content-Type": "application/json",
                                },
                                body: JSON.stringify({
                                    pageIndex: 0,
                                    zoomLevel: 11,
                                    latitude: marker.Position.Lat,
                                    longitude: marker.Position.Lng,
                                    geoHash: "f25ds",
                                }),
                            }
                        );

                        if (!res.ok) {
                            console.error(
                                "Error fetching marker info:",
                                res.status,
                                res.statusText
                            );
                            continue;
                        }

                        const data = await res.json();
                        console.log(`{ HTML: ${data.d.Result.Html} },`);
                    } catch (error) {
                        console.error(
                            "Error fetching marker info:",
                            error.message
                        );
                    }

                    // Delay for 1 second
                    await delay(1000);
                }

                console.log("]");
            } catch (error) {
                console.error("Error fetching markers:", error.message);
            }

            // Delay for 1 second before processing the next item
            await delay(1000);
        }
    });
})();
