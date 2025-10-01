 /** @type {import('tailwindcss').Config} */
 module.exports = {
 // NOTE Update this to include the paths to all files that contain Nativ
 ewind classes.
 content: ["./app/**/*.{js,jsx,ts,tsx}"],
 presets: [require("nativewind/preset")],
 theme: {
 extend: {},
 },
 plugins: [],
 }
 