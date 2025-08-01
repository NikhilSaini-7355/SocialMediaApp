import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.jsx'
import {ChakraProvider, extendTheme} from '@chakra-ui/react'
import { mode } from '@chakra-ui/theme-tools'
import { ColorModeScript } from '@chakra-ui/react'
import { BrowserRouter } from 'react-router-dom'

const styles = {
  global : (props)=> ({
    body : {
      color : mode('gray.800','whiteAlpha.900')(props),
      bg : mode('gray.100','#101010')(props)
    }
  })
}

const config = {
  initialColorMode : "dark",
  useSystemColorMode : true
}

const colors = {
  gray : {
    light : "#616161",
    dark : "#1e1e1e"
  }
};

const theme = extendTheme({
   config,
   styles,
   colors
});

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <BrowserRouter>
    <ChakraProvider theme={theme}>
    <ColorModeScript initialColorMode={theme.config.initialColorMode} />
    <App />
    </ChakraProvider>
    </BrowserRouter>
  </StrictMode>,
)
