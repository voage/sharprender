.PHONY: dev api web

dev:
	# Run both targets in parallel if you want: make -j2 dev
	make api & \
	make web

api:
	cd cmd/api && air

web:
	cd web && pnpm install && pnpm dev
