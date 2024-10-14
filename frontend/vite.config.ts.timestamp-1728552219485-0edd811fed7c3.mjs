// vite.config.ts
import { defineConfig } from "file:///C:/Users/geff/Desktop/new/galaxy3/frontend/node_modules/.pnpm/vite@5.4.1_@types+node@20.14.15_less@4.2.0/node_modules/vite/dist/node/index.js";
import vue from "file:///C:/Users/geff/Desktop/new/galaxy3/frontend/node_modules/.pnpm/@vitejs+plugin-vue@5.1.2_vite@5.4.1_@types+node@20.14.15_less@4.2.0__vue@3.4.38_typescript@5.4.5_/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import { fileURLToPath, URL } from "node:url";
import Components from "file:///C:/Users/geff/Desktop/new/galaxy3/frontend/node_modules/.pnpm/unplugin-vue-components@0.26.0_@babel+parser@7.25.3_rollup@4.20.0_vue@3.4.38_typescript@5.4.5_/node_modules/unplugin-vue-components/dist/vite.js";
import { AntDesignVueResolver } from "file:///C:/Users/geff/Desktop/new/galaxy3/frontend/node_modules/.pnpm/unplugin-vue-components@0.26.0_@babel+parser@7.25.3_rollup@4.20.0_vue@3.4.38_typescript@5.4.5_/node_modules/unplugin-vue-components/dist/resolvers.js";
var __vite_injected_original_import_meta_url = "file:///C:/Users/geff/Desktop/new/galaxy3/frontend/vite.config.ts";
var vite_config_default = defineConfig({
  plugins: [
    vue(),
    Components({
      resolvers: [
        AntDesignVueResolver({
          importStyle: false
          // css in js
        })
      ],
      types: [],
      dts: "src/components/components.d.ts",
      globs: ["src/components/*/index.vue"]
    })
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", __vite_injected_original_import_meta_url))
      // '@wails': fileURLToPath(new URL('./wailsjs', import.meta.url))
    }
  }
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcudHMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCJDOlxcXFxVc2Vyc1xcXFxnZWZmXFxcXERlc2t0b3BcXFxcbmV3XFxcXGdhbGF4eTNcXFxcZnJvbnRlbmRcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfZmlsZW5hbWUgPSBcIkM6XFxcXFVzZXJzXFxcXGdlZmZcXFxcRGVza3RvcFxcXFxuZXdcXFxcZ2FsYXh5M1xcXFxmcm9udGVuZFxcXFx2aXRlLmNvbmZpZy50c1wiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9pbXBvcnRfbWV0YV91cmwgPSBcImZpbGU6Ly8vQzovVXNlcnMvZ2VmZi9EZXNrdG9wL25ldy9nYWxheHkzL2Zyb250ZW5kL3ZpdGUuY29uZmlnLnRzXCI7aW1wb3J0IHsgZGVmaW5lQ29uZmlnIH0gZnJvbSAndml0ZSdcbmltcG9ydCB2dWUgZnJvbSAnQHZpdGVqcy9wbHVnaW4tdnVlJ1xuaW1wb3J0IHsgZmlsZVVSTFRvUGF0aCwgVVJMIH0gZnJvbSAnbm9kZTp1cmwnXG5pbXBvcnQgQ29tcG9uZW50cyBmcm9tICd1bnBsdWdpbi12dWUtY29tcG9uZW50cy92aXRlJ1xuaW1wb3J0IHsgQW50RGVzaWduVnVlUmVzb2x2ZXIgfSBmcm9tICd1bnBsdWdpbi12dWUtY29tcG9uZW50cy9yZXNvbHZlcnMnO1xuXG5cbi8vIGh0dHBzOi8vdml0ZWpzLmRldi9jb25maWcvXG5leHBvcnQgZGVmYXVsdCBkZWZpbmVDb25maWcoe1xuICBwbHVnaW5zOiBbXG4gICAgdnVlKCksXG4gICAgQ29tcG9uZW50cyh7XG4gICAgICByZXNvbHZlcnM6IFtcbiAgICAgICAgQW50RGVzaWduVnVlUmVzb2x2ZXIoe1xuICAgICAgICAgIGltcG9ydFN0eWxlOiBmYWxzZSwgLy8gY3NzIGluIGpzXG4gICAgICAgIH0pLFxuICAgICAgXSxcbiAgICAgIHR5cGVzOiBbXSxcbiAgICAgIGR0czogJ3NyYy9jb21wb25lbnRzL2NvbXBvbmVudHMuZC50cycsXG4gICAgICBnbG9iczogWydzcmMvY29tcG9uZW50cy8qL2luZGV4LnZ1ZSddXG4gICAgfSlcbiAgXSxcbiAgcmVzb2x2ZToge1xuICAgIGFsaWFzOiB7XG4gICAgICAnQCc6IGZpbGVVUkxUb1BhdGgobmV3IFVSTCgnLi9zcmMnLCBpbXBvcnQubWV0YS51cmwpKSxcbiAgICAgIC8vICdAd2FpbHMnOiBmaWxlVVJMVG9QYXRoKG5ldyBVUkwoJy4vd2FpbHNqcycsIGltcG9ydC5tZXRhLnVybCkpXG4gICAgfVxuICB9XG59KVxuIl0sCiAgIm1hcHBpbmdzIjogIjtBQUE4VCxTQUFTLG9CQUFvQjtBQUMzVixPQUFPLFNBQVM7QUFDaEIsU0FBUyxlQUFlLFdBQVc7QUFDbkMsT0FBTyxnQkFBZ0I7QUFDdkIsU0FBUyw0QkFBNEI7QUFKb0ssSUFBTSwyQ0FBMkM7QUFRMVAsSUFBTyxzQkFBUSxhQUFhO0FBQUEsRUFDMUIsU0FBUztBQUFBLElBQ1AsSUFBSTtBQUFBLElBQ0osV0FBVztBQUFBLE1BQ1QsV0FBVztBQUFBLFFBQ1QscUJBQXFCO0FBQUEsVUFDbkIsYUFBYTtBQUFBO0FBQUEsUUFDZixDQUFDO0FBQUEsTUFDSDtBQUFBLE1BQ0EsT0FBTyxDQUFDO0FBQUEsTUFDUixLQUFLO0FBQUEsTUFDTCxPQUFPLENBQUMsNEJBQTRCO0FBQUEsSUFDdEMsQ0FBQztBQUFBLEVBQ0g7QUFBQSxFQUNBLFNBQVM7QUFBQSxJQUNQLE9BQU87QUFBQSxNQUNMLEtBQUssY0FBYyxJQUFJLElBQUksU0FBUyx3Q0FBZSxDQUFDO0FBQUE7QUFBQSxJQUV0RDtBQUFBLEVBQ0Y7QUFDRixDQUFDOyIsCiAgIm5hbWVzIjogW10KfQo=
