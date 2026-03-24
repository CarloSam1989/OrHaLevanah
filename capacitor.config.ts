import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'com.moventanio.orhalevanah',
  appName: 'Or HaLevanah',
  webDir: 'mobile_web',
  bundledWebRuntime: false,
  server: {
    cleartext: true
  }
};

export default config;