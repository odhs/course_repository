# Ask Me Anything builded in Go + React

## Steps in Development

### Vite

```sh
npm create vite@latest
```

Select:

- web
- React
- Typescript

Install Dependences

```sh
npm i
```

### CSS

Add the CSS framework

```sh
npm i tailwindcss postcss autoprefixer -D
npx tailwindcss init -p
```

### Router

Install the router

```sh
npm i react-router-dom
```

### Icons

Add icons package

```sh
npm i lucide-react
```

### Migrate to React 19

To get the form parameters with new feature of React 19

```sh
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

### Install Sonner to show a toaster

To show messages to the user, -f because we are using React 19 and this version is in release candidate yet. The function was added to the Sharing Button in Room page.

```sh
npm i sonner -f
```

### Install QueryClient

```sh
npm i @tanstack/react-query -f
```

## Run

### To start the application in _development mode_

```sh
npm run dev
```

### To start the application to preview in entire network

```sh
npm run expose
```

PS: I added the line `"expose": "npm run dev -- --host"` into the term `scripts` in the file `package.json` to test the application on my network, this line permits that other computers or mobiles can access the website in a same network, for example, a Wifi network.
