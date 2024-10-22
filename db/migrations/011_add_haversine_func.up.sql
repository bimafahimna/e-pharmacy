CREATE OR REPLACE FUNCTION haversine_distance(
    lat1 DECIMAL,
    lon1 DECIMAL,
    lat2 DECIMAL,
    lon2 DECIMAL
) RETURNS DECIMAL AS $$
DECLARE
    r DECIMAL := 6371; -- Radius of the Earth in kilometers
    lat1_rad DECIMAL;
    lat2_rad DECIMAL;
    delta_lat DECIMAL;
    delta_lon DECIMAL;
    a DECIMAL;
    c DECIMAL;
BEGIN
    -- Convert degrees to radians
    lat1_rad := radians(lat1);
    lat2_rad := radians(lat2);
    delta_lat := radians(lat2 - lat1);
    delta_lon := radians(lon2 - lon1);

    -- Haversine formula
    a := sin(delta_lat / 2) * sin(delta_lat / 2) +
         cos(lat1_rad) * cos(lat2_rad) *
         sin(delta_lon / 2) * sin(delta_lon / 2);
    c := 2 * atan2(sqrt(a), sqrt(1 - a));

    -- Distance in kilometers
    RETURN r * c;
END;
$$ LANGUAGE plpgsql;