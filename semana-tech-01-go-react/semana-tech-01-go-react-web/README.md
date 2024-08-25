# Iniciando a aplicação

```bs
npm create vite@latest
```

Selecione:

- web
- REACT
- Typescript

Entra na pasta e roda no VSCode

```bs
cd web
code .
```

Instalando as dependências

```bs
npm i
```

Adicionando o CSS

```bs
npm i tailwindcss postcss autoprefixer -D
npx tailwindcss init -p
```

Instalando o roteador

```bs
npm i react-router-dom
```

Adicionando ícone

```bs
npm i lucide-react
```

Migrando para REACT 19

```bs
npm install --save-exact react@rc react-dom@rc
```

1 Substituir no package.json:

```json
  "devDependencies": {
    "@eslint/js": "^9.9.0",
    "@types/react": "npm:types-react@rc",
    "@types/react-dom": "npm:types-react-dom@rc",
  }
```

2 Adicionar no package.json

```json
  "overrides": {
    "@types/react": "npm:types-react@rc",
    "@types/react-dom": "npm:types-react-dom@rc"
  }
```

3 Deletar:

- pasta: node_modules
- arquivo: package-lock.json

4: Executar `npm i -f` para instalar tudo

## Para iniciar a aplicação use

```bs
npm run dev
```
