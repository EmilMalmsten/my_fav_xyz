import './App.css'
import 'bootstrap/dist/css/bootstrap.min.css';
import { Route, Routes } from 'react-router-dom';
import MainNavbar from './components/MainNavbar';
import Home from './Home';
import Register from './Register';
import Login from './Login';
import Toplist from './Toplist';
import ToplistsByCategory from './components/ToplistsByCategory';

function App() {

  return (
    <>
      <MainNavbar/>
      <Routes>
        <Route path='/' element={<Home/>}></Route>
        <Route path='/register' element={<Register/>}></Route>
        <Route path='/login' element={<Login/>}></Route>
        <Route path='/toplists/:id' element={<Toplist/>}></Route>
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
