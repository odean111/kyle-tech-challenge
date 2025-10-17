import React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import axios from 'axios';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';

const companySchema = z.object({
  jurisdiction: z.enum(['UK', 'Singapore', 'Caymens']),
  company_name: z.string().min(1, 'Company name is required').max(255),
  company_address: z.string().min(1, 'Company address is required'),
  nature_of_business: z.string().optional(),
  number_of_directors: z.number().min(1).max(100).optional(),
  number_of_shareholders: z.number().min(1).max(1000).optional(),
  sec_code: z.string().optional(),
});

type CompanyFormData = z.infer<typeof companySchema>;

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

const CreateCompany: React.FC = () => {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const isEditing = Boolean(id);

  const form = useForm<CompanyFormData>({
    resolver: zodResolver(companySchema),
    defaultValues: {
      jurisdiction: 'UK',
      company_name: '',
      company_address: '',
      nature_of_business: '',
      number_of_directors: undefined,
      number_of_shareholders: undefined,
      sec_code: '',
    },
  });

  const { handleSubmit, control, setValue, formState: { isSubmitting } } = form;

  const fetchCompany = React.useCallback(async (companyId: string) => {
    try {
      const response = await axios.get<Company>(`http://localhost:8080/api/v1/companies/${companyId}`);
      const company = response.data;
      
      setValue('jurisdiction', company.jurisdiction as CompanyFormData['jurisdiction']);
      setValue('company_name', company.company_name);
      setValue('company_address', company.company_address);
      setValue('nature_of_business', company.nature_of_business || '');
      setValue('number_of_directors', company.number_of_directors);
      setValue('number_of_shareholders', company.number_of_shareholders);
      setValue('sec_code', company.sec_code || '');
    } catch (error) {
      console.error('Error fetching company:', error);
      alert('Failed to fetch company data');
    }
  }, [setValue]);

  React.useEffect(() => {
    if (isEditing && id) {
      fetchCompany(id);
    }
  }, [isEditing, id, fetchCompany]);

  const onSubmit = async (data: CompanyFormData) => {
    try {
      if (isEditing && id) {
        await axios.put(`http://localhost:8080/api/v1/companies/${id}`, data);
      } else {
        await axios.post('http://localhost:8080/api/v1/companies', data);
      }
      navigate('/');
    } catch (error) {
      console.error('Error saving company:', error);
      alert('Failed to save company');
    }
  };

  return (
    <div className="container mx-auto p-6 max-w-2xl">
      <Card>
        <CardHeader>
          <CardTitle className="text-2xl font-bold">
            {isEditing ? 'Edit Company' : 'Create New Company'}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              <FormField
                control={control}
                name="company_name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Company Name</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter company name" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={control}
                name="company_address"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Company Address</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter company address" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={control}
                name="jurisdiction"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Jurisdiction</FormLabel>
                    <Select onValueChange={field.onChange} defaultValue={field.value}>
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select jurisdiction" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="UK">UK</SelectItem>
                        <SelectItem value="Singapore">Singapore</SelectItem>
                        <SelectItem value="Caymens">Caymens</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={control}
                name="nature_of_business"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Nature of Business (Optional)</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter nature of business" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={control}
                name="number_of_directors"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Number of Directors (Optional)</FormLabel>
                    <FormControl>
                      <Input 
                        type="number" 
                        placeholder="Enter number of directors"
                        {...field}
                        onChange={(e) => field.onChange(e.target.value ? parseInt(e.target.value) : undefined)}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={control}
                name="number_of_shareholders"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Number of Shareholders (Optional)</FormLabel>
                    <FormControl>
                      <Input 
                        type="number" 
                        placeholder="Enter number of shareholders"
                        {...field}
                        onChange={(e) => field.onChange(e.target.value ? parseInt(e.target.value) : undefined)}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={control}
                name="sec_code"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>SEC Code (Optional)</FormLabel>
                    <FormControl>
                      <Input placeholder="Enter SEC code" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="flex gap-4">
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting ? 'Saving...' : (isEditing ? 'Update Company' : 'Create Company')}
                </Button>
                <Button type="button" variant="outline" onClick={() => navigate('/')}>
                  Cancel
                </Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
};

export default CreateCompany;