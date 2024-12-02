const output = document.getElementById('output');
const orderWidget = document.getElementById('order-widget');
const orderDetailsDiv = document.getElementById('order-details');
const toggleButton = document.getElementById('close-button');

document.getElementById('getButton').addEventListener('click', () => {
    const orderId = document.getElementById('orderInput').value;
    output.textContent = 'Загрузка...';

    fetch(`http://localhost:8080/WB/get?order_uuid=${encodeURIComponent(orderId)}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`Ошибка: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            output.textContent = JSON.stringify(data, null, 2);
            displayOrderDetails(data);
        })
        .catch(error => {
            output.textContent = `Ошибка: ${error.message}`;
        });
});

function displayOrderDetails(order) {
    orderWidget.style.display = 'block';
    let htmlContent = `
        <p><strong>Order UID:</strong> ${order.order_uid}</p>
        <p><strong>Track Number:</strong> ${order.track_number}</p>
        <p><strong>Customer ID:</strong> ${order.customer_id}</p>
        <p><strong>Delivery Service:</strong> ${order.delivery_service}</p>
        <p><strong>Date Created:</strong> ${new Date(order.date_created).toLocaleString()}</p>
        <h4>Items</h4>
    `;

    order.items.forEach(item => {
        htmlContent += `
            <div>
                <p><strong>Product Name:</strong> ${item.name}</p>
                <p><strong>Price:</strong> ${(item.price / 100).toFixed(2)}</p>
                <p><strong>Quantity:</strong> ${item.sale}</p>
            </div>
        `;
    });

    orderDetailsDiv.innerHTML = htmlContent;
}

toggleButton.addEventListener('click', () => {
    orderWidget.style.display = 'none';
});
