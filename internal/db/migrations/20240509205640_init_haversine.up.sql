CREATE OR REPLACE FUNCTION haversine(
	lat1 FLOAT8, lon1 FLOAT8, lat2 FLOAT8, lon2 FLOAT8
)
RETURNS FLOAT8 AS $$
DECLARE
	lat1_rad FLOAT8;
	lat2_rad FLOAT8;
	delta_lat FLOAT8;
	delta_lon FLOAT8;
	a FLOAT8;
	c FLOAT8;
BEGIN
	lat1_rad = RADIANS(lat1);
	lat2_rad = RADIANS(lat2);
	delta_lat = RADIANS(lat2 - lat1);
	delta_lon = RADIANS(lon2 - lon1);
	a = SIN(delta_lat / 2) * SIN(delta_lat / 2) + COS(lat1_rad) * COS(lat2_rad) * SIN(delta_lon / 2) * SIN(delta_lon / 2);
	c = 2 * ATAN2(SQRT(a), SQRT(1 - a));
	RETURN 6371 * c;
END;
$$ LANGUAGE plpgsql;