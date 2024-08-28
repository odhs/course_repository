# Ask Me Anything builded in Go + React

## Steps in Development

### Vite

```bs
npm create vite@latest
```

Select:

- web
- REACT
- Typescript

Install Dependences

```bs
npm i
```

### CSS

Add the CSS framework

```bs
npm i tailwindcss postcss autoprefixer -D
npx tailwindcss init -p
```

### Router

Install the router

```bs
npm i react-router-dom
```

### Icons

Add icons package

```bs
npm i lucide-react
```

### Migrating to React 19

To get the form parameters with new feature of React 19

```bs
npm install --save-exact react@rc react-dom@rc
```

1 Replace in package.json:

```json
  "devDependencies": {
    "@eslint/js": "^9.9.0",
    "@types/react": "npm:types-react@rc",
    "@types/react-dom": "npm:types-react-dom@rc",
  }
```

2 Add to package.json

```json
  "overrides": {
    "@types/react": "npm:types-react@rc",
    "@types/react-dom": "npm:types-react-dom@rc"
  }
```

3 Delete:

- directory: node_modules
- file: package-lock.json

4: Execute `npm i -f` to install everything

### Instaling package Sonner to show a toaster

To show messages to the user, -f because we are using React 19 and this version is in release candidate yet. The function was added to the Sharing Button in Room page.

```bs
npm i sonner -f
```

## To start the application use

```bs
npm run dev
```
