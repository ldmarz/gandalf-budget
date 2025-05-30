const API_BASE_URL = '/api/v1';

interface ApiErrorResponse {
  error: string;
  details?: any;
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let errorData: ApiErrorResponse = { error: `HTTP error! status: ${response.status}` };
    try {
      const parsedError = await response.json();
      if (parsedError && parsedError.error) {
        errorData = parsedError as ApiErrorResponse;
      }
    } catch (e) {
      console.error("Could not parse error response as JSON:", e);
    }
    console.error('API call failed:', errorData);
    throw new Error(errorData.error);
  }
  const contentType = response.headers.get("content-type");
  if (contentType && contentType.indexOf("application/json") !== -1) {
    return response.json() as Promise<T>;
  } else if (response.status === 204) {
    return Promise.resolve(null as T);
  }
  return response.text() as Promise<any>;
}

export async function get<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
  });
  return handleResponse<T>(response);
}

export async function post<T, U>(path: string, body: U): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
    body: JSON.stringify(body),
  });
  return handleResponse<T>(response);
}

export async function put<T, U>(path: string, body: U): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
    body: JSON.stringify(body),
  });
  return handleResponse<T>(response);
}

export async function del<T>(path: string): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: 'DELETE',
    headers: {
      'Accept': 'application/json',
    },
  });
  return handleResponse<T>(response);
}

export interface BudgetLine {
  id: number;
  month_id: number;
  category_id: number;
  label: string;
  expected: number;
  category_name?: string;
  category_color?: string;
  actual_amount?: number;
  actual_id?: number;
}

export interface ActualLine {
  id: number;
  budget_line_id: number;
  actual: number;
}

// API functions for Budget Lines
export async function createBudgetLine(data: { month_id: number; category_id: number; label: string; expected: number }): Promise<BudgetLine> {
  return post<BudgetLine, typeof data>('/budget-lines', data);
}

export async function getBudgetLinesByMonth(monthId: number): Promise<BudgetLine[]> {
  return get<BudgetLine[]>(`/budget-lines?month_id=${monthId}`);
}

export async function updateBudgetLine(id: number, data: { label?: string; expected?: number }): Promise<BudgetLine> {
  return put<BudgetLine, typeof data>(`/budget-lines/${id}`, data);
}

export async function deleteBudgetLine(id: number): Promise<void> {
  return del<void>(`/budget-lines/${id}`);
}

export async function updateActualLine(id: number, data: { actual: number }): Promise<ActualLine> {
  return put<ActualLine, typeof data>(`/actual-lines/${id}`, data);
}

export interface BudgetLineWithActual {
  id: number;
  month_id: number;
  category_id: number;
  category_name: string;
  category_color: string;
  label: string;
  expected_amount: number;
  actual_amount: number;
}

export interface BoardDataPayload {
  month_id: number;
  year: number;
  month_name: string;
  budget_lines: BudgetLineWithActual[];
}

export async function getBoardData(monthId: string | number): Promise<BoardDataPayload> {
  return get<BoardDataPayload>(`/board-data/${monthId}`);
}

interface FinalizeMonthResponse {
  message: string;
  new_month_id: number;
}

export async function finalizeMonth(monthId: string | number): Promise<FinalizeMonthResponse> {
  return put<FinalizeMonthResponse, null>(`/months/${monthId}/finalize`, null);
}

export interface Category {
  id: number;
  name: string;
  color: string;
}

export async function getAllCategories(): Promise<Category[]> {
  return get<Category[]>('/categories');
}

export async function createCategory(data: { name: string; color: string }): Promise<Category> {
  return post<Category, typeof data>('/categories', data);
}

export async function updateCategory(id: number, data: { name?: string; color?: string }): Promise<Category> {
  return put<Category, typeof data>(`/categories/${id}`, data);
}

export async function deleteCategory(id: number): Promise<void> {
  return del<void>(`/categories/${id}`);
}

export interface BudgetLineDetail {
  budget_line_id: number;
  label: string;
  expected_amount: number;
  actual_amount: number;
  difference: number;
}

export interface CategorySummary {
  category_id: number;
  category_name: string;
  category_color: string;
  total_expected: number;
  total_actual: number;
  difference: number;
  budget_lines: BudgetLineDetail[];
}

export interface DashboardPayload {
  month_id: number;
  year: number;
  month: string;
  total_expected: number;
  total_actual: number;
  total_difference: number;
  category_summaries: CategorySummary[];
}

export async function getDashboardData(monthId: number | string): Promise<DashboardPayload> {
  return get<DashboardPayload>(`/dashboard?month_id=${monthId}`);
}

export interface AnnualSnapMeta {
  id: number;
  month_id: number;
  year: number;
  month: string;
  snap_created_at: string;
}

export async function getAnnualSnapshots(year: number): Promise<AnnualSnapMeta[]> {
  return get<AnnualSnapMeta[]>(`/reports/annual?year=${year}`);
}

export async function getSnapshotDetail(snapId: number): Promise<DashboardPayload> {
  return get<DashboardPayload>(`/reports/snapshots/${snapId}`);
}
