import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import CompanyList from './pages/CompanyList';
import CreateCompany from './pages/CreateCompany';

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-background">
        <Routes>
          <Route path="/" element={<CompanyList />} />
          <Route path="/create" element={<CreateCompany />} />
          <Route path="/edit/:id" element={<CreateCompany />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
