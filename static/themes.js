'use strict';

const THEMES = {
  dark: {
    name: 'Dark',
    vars: {
      '--bg':              '#000000',
      '--bg-card':         '#000000',
      '--bg-hover':        '#111111',
      '--border':          '#FFFFFF',
      '--border-dim':      '#888888',
      '--green':           '#00FF41',
      '--amber':           '#FFC800',
      '--red':             '#FF5555',
      '--text':            '#FFFFFF',
      '--text-dim':        '#E0E0E0',
    },
  },
  light: {
    name: 'Light',
    vars: {
      '--bg':              '#FFFFFF',
      '--bg-card':         '#F5F5F5',
      '--bg-hover':        '#EEEEEE',
      '--border':          '#222222',
      '--border-dim':      '#666666',
      '--green':           '#00AA30',
      '--amber':           '#CC9900',
      '--red':             '#DD3333',
      '--text':            '#111111',
      '--text-dim':        '#555555',
    },
  },
  highContrast: {
    name: 'High Contrast',
    vars: {
      '--bg':              '#000000',
      '--bg-card':         '#000000',
      '--bg-hover':        '#222222',
      '--border':          '#FFFF00',
      '--border-dim':      '#FFFF00',
      '--green':           '#00FF00',
      '--amber':           '#FFFF00',
      '--red':             '#FF0000',
      '--text':            '#FFFFFF',
      '--text-dim':        '#CCCCCC',
    },
  },
};

const STORAGE_KEY = 'dashboard-theme';

function getStoredTheme() {
  try {
    return localStorage.getItem(STORAGE_KEY);
  } catch {
    return null;
  }
}

function storeTheme(themeName) {
  try {
    localStorage.setItem(STORAGE_KEY, themeName);
  } catch {
  }
}

function applyTheme(themeName) {
  const theme = THEMES[themeName];
  if (!theme) return;
  const root = document.documentElement;
  root.setAttribute('data-theme', themeName);
  for (const [prop, value] of Object.entries(theme.vars)) {
    root.style.setProperty(prop, value);
  }
  storeTheme(themeName);
}

function loadTheme() {
  const stored = getStoredTheme();
  const theme = stored && THEMES[stored] ? stored : 'dark';
  applyTheme(theme);
  return theme;
}

function cycleTheme() {
  const themeNames = Object.keys(THEMES);
  const stored = getStoredTheme();
  let idx = themeNames.indexOf(stored);
  if (idx < 0) idx = 0;
  idx = (idx + 1) % themeNames.length;
  applyTheme(themeNames[idx]);
  return themeNames[idx];
}