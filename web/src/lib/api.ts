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
