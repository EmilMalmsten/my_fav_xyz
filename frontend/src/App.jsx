import './App.css'
import 'bootstrap/dist/css/bootstrap.min.css';
import { Route, Routes } from 'react-router-dom';
import MainNavbar from './components/MainNavbar';
import Home from './pages/Home';
import Register from './pages/Register';
import Login from './pages/Login';
import Toplist from './pages/ViewToplist';
import CreateToplist from './pages/CreateToplist';
import ToplistsByCategory from './components/ToplistsByCategory';
import TokenManager from './services/TokenManager';

function App() {

  return (
    <>
      <TokenManager/>
      <MainNavbar/>
      <Routes>
        <Route path='/' element={<Home/>}></Route>
        <Route path='/register' element={<Register/>}></Route>
        <Route path='/login' element={<Login/>}></Route>
        <Route path='/toplists/:id' element={<Toplist/>}></Route>
        <Route path='/toplists/create' element={<CreateToplist/>}></Route>
        <Route path='/toplists/recent' element={
          <ToplistsByCategory title="Most recent toplists" endpoint="/toplists/recent"/>
        }>
        </Route>
        <Route path='/toplists/popular' element={
          <ToplistsByCategory title="Most popular toplists" endpoint="/toplists/popular"/>
        }>
        </Route>
      </Routes>
    </>
  )
}

export default App
