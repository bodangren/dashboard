---
version: 1.0.0
name: Dashboard Design System
colors:
  background: "#000000"
  surface: "#000000"
  primary: "#00FF41"
  secondary: "#FFC800"
  error: "#FF5555"
  text: "#FFFFFF"
  text-dim: "#E0E0E0"
  border: "#FFFFFF"
  border-dim: "#888888"
typography:
  headline-lg:
    fontFamily: '"JetBrains Mono", "Fira Code", "Cascadia Code", monospace'
    fontSize: 24px
    fontWeight: 700
    lineHeight: 1.2
  body-md:
    fontFamily: '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif'
    fontSize: 16px
    fontWeight: 400
    lineHeight: 1.6
  code-md:
    fontFamily: '"JetBrains Mono", "Fira Code", "Cascadia Code", monospace'
    fontSize: 14px
    fontWeight: 400
    lineHeight: 1.6
spacing:
  base: 16px
rounded:
  none: 0px
---

# Design

This document defines the design system for the Dashboard project, characterized by a high-contrast, terminal-inspired aesthetic. It prioritizes data density and readability.

## Colors

The color palette is designed for high visibility in dark environments, utilizing a limited set of vibrant accent colors against a pure black background.

- **Background (#000000):** The primary background color for the entire application.
- **Surface (#000000):** Card and container backgrounds, identical to the background to maintain a "flat" terminal feel.
- **Primary (#00FF41):** "Matrix Green," used for success states, active highlights, primary buttons, and commit hashes.
- **Secondary (#FFC800):** Amber/Gold, used for warnings, meta-information in diffs, and secondary highlights.
- **Error (#FF5555):** Soft Red, used for error states and destructive actions.
- **Text (#FFFFFF):** Pure white for high-readability primary text.
- **Text Dim (#E0E0E0):** Off-white/gray for secondary information and metadata.
- **Border (#FFFFFF):** Primary borders and separators.
- **Border Dim (#888888):** Subtle separators for list items and secondary elements.

## Typography

The typography system uses a mix of system-native sans-serif for UI elements and monospace for data-heavy content.

- **Headline LG:** Used for page headers. Monospace, bold, 24px.
- **Body MD:** Used for standard UI labels and text. System sans-serif, 16px.
- **Code MD:** Used for commit hashes, logs, and metadata. Monospace, 14px.

## Layout

The application uses a mobile-first, responsive layout based on a 16px grid.

- **Spacing Base:** 16px.
- **Gap:** 1rem (16px) on mobile, 1.5rem (24px) on desktop.
- **Padding:** Containers use 1rem to 2rem of padding depending on screen size.
- **Max Width:** The main content area is capped at 1600px.
- **Grid:** Project cards use a responsive grid that shifts from 1 to 2 columns at the 1025px breakpoint.

## Elevation & Depth

The design is strictly flat to reinforce the terminal aesthetic.

- **Elevation:** None. No box-shadows or drop-shadows are used.
- **Depth:** Indicated solely through borders and color shifts on hover.

## Shapes

The design avoids rounded corners to maintain a strict, industrial aesthetic.

- **None (0px):** All buttons, cards, and input fields use a border-radius of 0.
- **Exception:** Spinners use a 50% radius for circular motion.

## Components

### Buttons
Buttons are outlined with a 2px border. The primary button uses the Primary Green color, while secondary buttons use the Border color.
- **Hover:** Background fills with the border color, and text flips to the background color.

### Cards
Project and agent cards use a 1px or 2px white border with no rounded corners.

### Nav Tabs
Navigation tabs use a heavy 3px bottom border for the active state and a dim border for inactive states.

### Diffs
Diffs use high-contrast coloring for additions (Primary Green) and removals (Error Red), with amber highlights for metadata.

### Iconography
Icons are SVG-based, using a 2px stroke width with square caps and miter joins.
- **Default Size:** 20px - 24px.

## Do's and Don'ts

### Do
- Use pure black (#000000) for backgrounds to ensure maximum contrast.
- Use monospace fonts for all data-driven content (hashes, logs, paths).
- Maintain sharp 0px corners on all interactive elements.
- Use Primary Green (#00FF41) for positive actions and success states.

### Don't
- Do not use gradients or shadows; the design is strictly flat.
- Do not use rounded corners (except for spinners).
- Do not introduce low-contrast grays for primary text.
- Do not use transparency unless indicating a disabled state.
