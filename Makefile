.PHONY: dev build install uninstall clean

WAILS := $(shell command -v wails 2>/dev/null || echo "$(shell go env GOPATH)/bin/wails")
PREFIX := $(HOME)/.local

dev:
	$(WAILS) dev

build:
	$(WAILS) build

install: build
	mkdir -p $(PREFIX)/bin $(PREFIX)/share/icons $(PREFIX)/share/applications
	cp build/bin/sushi $(PREFIX)/bin/sushi
	cp build/appicon.png $(PREFIX)/share/icons/sushi.png
	@echo '[Desktop Entry]' > $(PREFIX)/share/applications/sushi.desktop
	@echo 'Name=Sushi' >> $(PREFIX)/share/applications/sushi.desktop
	@echo 'Comment=3D Asset Manager' >> $(PREFIX)/share/applications/sushi.desktop
	@echo 'Exec=sushi' >> $(PREFIX)/share/applications/sushi.desktop
	@echo 'Icon=sushi' >> $(PREFIX)/share/applications/sushi.desktop
	@echo 'Type=Application' >> $(PREFIX)/share/applications/sushi.desktop
	@echo 'Categories=Graphics;3DGraphics;Utility;' >> $(PREFIX)/share/applications/sushi.desktop
	@echo 'Terminal=false' >> $(PREFIX)/share/applications/sushi.desktop
	@echo "✅ Installed sushi to $(PREFIX)/bin/sushi"

uninstall:
	rm -f $(PREFIX)/bin/sushi
	rm -f $(PREFIX)/share/icons/sushi.png
	rm -f $(PREFIX)/share/applications/sushi.desktop
	@echo "🗑  Uninstalled sushi"

update: install
	@echo "🔄 Updated sushi"

clean:
	rm -rf build/bin/sushi frontend/dist
