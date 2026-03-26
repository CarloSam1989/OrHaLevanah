import type { CapacitorConfig } from '@capacitor/cli';

const config: CapacitorConfig = {
  appId: 'com.orhalevanah.app',
  appName: 'Or HaLevanah',
  webDir: 'www',

  plugins: {
    SplashScreen: {
      launchShowDuration: 3000, // tiempo en ms (3 segundos)
      launchAutoHide: true,     // se oculta automáticamente
      backgroundColor: "#000000", // color de fondo (puedes usar tu color del logo)
      androidSplashResourceName: "splash", // nombre del archivo
      androidScaleType: "CENTER_CROP",
      showSpinner: false
    }
  }
};

export default config;