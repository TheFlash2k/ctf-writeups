import { MantineProvider } from '@mantine/core'
import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import { BrowserRouter } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from 'react-query'


const queryClient = new QueryClient()


ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
     <MantineProvider withGlobalStyles withNormalizeCSS>
      <BrowserRouter>
      <App />
      </BrowserRouter>
    </MantineProvider>
    </QueryClientProvider>
  </React.StrictMode>,
)