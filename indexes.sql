CREATE INDEX new_york_geom_idx
ON public.new_york
USING GIST (geom);
