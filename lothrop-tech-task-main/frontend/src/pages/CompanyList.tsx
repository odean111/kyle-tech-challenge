import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import axios from 'axios';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

interface Company {
  id: string;
  jurisdiction: string;
  company_name: string;
  company_address: string;
  nature_of_business?: string;
  number_of_directors?: number;
  number_of_shareholders?: number;
  sec_code?: string;
  date_created: string;
  date_updated: string;
}

interface CompaniesResponse {
  companies: Company[];
  limit: number;
  offset: number;
  total: number;
}

const CompanyList: React.FC = () => {
  const [companies, setCompanies] = useState<Company[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchCompanies();
  }, []);

  const fetchCompanies = async () => {
    try {
      setLoading(true);
      const response = await axios.get<CompaniesResponse>('http://localhost:8080/api/v1/companies');
      setCompanies(response.data.companies);
      setError(null);
    } catch (err) {
      setError('Failed to fetch companies');
      console.error('Error fetching companies:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!window.confirm('Are you sure you want to delete this company?')) {
      return;
    }

    try {
      await axios.delete(`http://localhost:8080/api/v1/companies/${id}`);
      await fetchCompanies(); // Refresh the list
    } catch (err) {
      setError('Failed to delete company');
      console.error('Error deleting company:', err);
    }
  };

  if (loading) {
    return (
      <div className="container mx-auto p-6">
        <Card>
          <CardContent className="flex items-center justify-center p-6">
            <div className="text-center">Loading companies...</div>
          </CardContent>
        </Card>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto p-6">
        <Card>
          <CardContent className="p-6">
            <div className="text-center text-red-600">{error}</div>
            <div className="text-center mt-4">
              <Button onClick={fetchCompanies}>Try Again</Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-6">
      <Card>
        <CardHeader>
          <div className="flex justify-between items-center">
            <CardTitle className="text-2xl font-bold">Companies</CardTitle>
            <Button asChild>
              <Link to="/create">Add New Company</Link>
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          {companies.length === 0 ? (
            <div className="text-center py-8">
              <p className="text-gray-500 mb-4">No companies found</p>
              <Button asChild>
                <Link to="/create">Create your first company</Link>
              </Button>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Company Name</TableHead>
                  <TableHead>Address</TableHead>
                  <TableHead>Jurisdiction</TableHead>
                  <TableHead>Business Nature</TableHead>
                  <TableHead>Directors</TableHead>
                  <TableHead>Shareholders</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {companies.map((company) => (
                  <TableRow key={company.id}>
                    <TableCell className="font-medium">{company.company_name}</TableCell>
                    <TableCell>{company.company_address}</TableCell>
                    <TableCell>
                      <Badge variant="outline">{company.jurisdiction}</Badge>
                    </TableCell>
                    <TableCell>{company.nature_of_business || 'N/A'}</TableCell>
                    <TableCell>{company.number_of_directors || 'N/A'}</TableCell>
                    <TableCell>{company.number_of_shareholders || 'N/A'}</TableCell>
                    <TableCell className="text-right">
                      <div className="flex gap-2 justify-end">
                        <Button variant="outline" size="sm" asChild>
                          <Link to={`/edit/${company.id}`}>Edit</Link>
                        </Button>
                        <Button 
                          variant="destructive" 
                          size="sm"
                          onClick={() => handleDelete(company.id)}
                        >
                          Delete
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default CompanyList;