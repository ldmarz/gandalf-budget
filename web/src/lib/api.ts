const API_BASE_URL = '/api/v1';

interface ApiErrorResponse {
  error: string;
  details?: any; // Or a more specific error details type
}

// Helper function to handle responses
async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let errorData: ApiErrorResponse = { error: `HTTP error! status: ${response.status}` };
    try {
      // Try to parse a JSON error response from the backend
      const parsedError = await response.json();
      if (parsedError && parsedError.error) {
        errorData = parsedError as ApiErrorResponse;
      }
    } catch (e) {
      // Could not parse JSON error, stick with the HTTP status
      console.error("Could not parse error response as JSON:", e);
    }
    console.error('API call failed:', errorData);
    throw new Error(errorData.error); // Throw an error that can be caught by the caller
  }
  // If response is OK, try to parse JSON, but handle cases with no content (e.g., 204)
  const contentType = response.headers.get("content-type");
  if (contentType && contentType.indexOf("application/json") !== -1) {
    return response.json() as Promise<T>;
  } else if (response.status === 204) { // No Content
    return Promise.resolve(null as T); // Or undefined, depending on how you want to handle it
  }
  // For non-JSON responses, you might want to return text or blob
  // For now, we'll assume JSON or no content for successful responses
  return response.text() as Promise<any>; // Fallback for unexpected content types
}

// GET request helper
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

// POST request helper
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

// PUT request helper
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

// DELETE request helper
export async function del<T>(path: string): Promise<T> { // 'delete' is a reserved keyword
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: 'DELETE',
    headers: {
      'Accept': 'application/json', // Expect JSON error response, but 204 on success
    },
  });
  return handleResponse<T>(response); // Will handle 204 No Content correctly
}

// Example usage (optional, for testing or demonstration):
/*
interface Category {
  id: number;
  name: string;
  color: string;
}

async function testApi() {
  try {
    const categories = await get<Category[]>('/categories');
    console.log('Categories:', categories);

    if (categories.length > 0) {
      const firstCatId = categories[0].id;
      const updatedCat = await put<Category, Partial<Category>>(`/categories/${firstCatId}`, { name: 'Updated Name' });
      console.log('Updated Category:', updatedCat);
    }

  } catch (error) {
    console.error('Error in testApi:', error);
  }
}
// testApi();
*/

// TypeScript types for BudgetLine and ActualLine
export interface BudgetLine {
  id: number;
  month_id: number;
  category_id: number;
  label: string;
  expected: number;
  // Optional fields if your API might join them, or if you plan to merge client-side
  category_name?: string; 
  category_color?: string;
  actual_amount?: number; // To store the actual amount from ActualLine
  actual_id?: number; // ID of the associated ActualLine record
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
  return del<void>(`/budget-lines/${id}`); // Expects 204 No Content, handleResponse handles this
}

// API function for Actual Lines
export async function updateActualLine(id: number, data: { actual: number }): Promise<ActualLine> {
  return put<ActualLine, typeof data>(`/actual-lines/${id}`, data);
}

// --- Board Data API (Updated) ---
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
  // categories?: Category[]; // This was optional, not adding for now as backend doesn't include it here
}

// API function for Board Data (Updated)
export async function getBoardData(monthId: string | number): Promise<BoardDataPayload> {
  // The PRD endpoint for board data is /board/:monthId, but current implementation uses /board-data/:monthId
  // Assuming /board-data/:monthId is correct as per existing code.
  // If it should be /api/v1/board/:monthId, this path needs to be changed.
  // Sticking to /board-data/ for now.
  return get<BoardDataPayload>(`/board-data/${monthId}`);
}

// API function for Finalizing Month
interface FinalizeMonthResponse {
  message: string; // e.g., "Month finalized successfully"
  new_month_id: number;
}

export async function finalizeMonth(monthId: string | number): Promise<FinalizeMonthResponse> {
  return put<FinalizeMonthResponse, null>(`/months/${monthId}/finalize`, null); // No body for this PUT request
}

// --- Existing Category types and functions for context ---
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

// --- Dashboard API ---

export interface BudgetLineDetail {
  budget_line_id: number; // Note: Go uses `json:"budget_line_id"` which becomes budget_line_id
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

// --- Reports API ---

export interface AnnualSnapMeta {
  id: number; // Changed from int64 in Go to number in TS
  month_id: number; // Changed from int64 in Go to number in TS
  year: number;
  month: string; 
  snap_created_at: string; // ISO date string (Go's time.Time marshals to this format)
}

export async function getAnnualSnapshots(year: number): Promise<AnnualSnapMeta[]> {
  return get<AnnualSnapMeta[]>(`/reports/annual?year=${year}`);
}

// The backend GetSnapshotDetail returns the raw JSON string of DashboardPayload.
// The `get` helper in api.ts already handles response.json(), so it should correctly parse this.
export async function getSnapshotDetail(snapId: number): Promise<DashboardPayload> {
  return get<DashboardPayload>(`/reports/snapshots/${snapId}`);
}
