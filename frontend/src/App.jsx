import './App.css'
import 'bootstrap/dist/css/bootstrap.min.css';
import { Route, Routes } from 'react-router-dom';
import MainNavbar from './components/main_navbar';
import Home from './Home';
import Register from './Register';
import Login from './Login';
import Toplist from './Toplist';

function App() {

  return (
    <>
      <MainNavbar/>
      <Routes>
        <Route path='/' element={<Home/>}></Route>
        <Route path='/register' element={<Register/>}></Route>
        <Route path='/login' element={<Login/>}></Route>
        <Route path='/toplists/:id' element={<Toplist/>}></Route>
      </Routes>
    </>
  )
}

export default App
