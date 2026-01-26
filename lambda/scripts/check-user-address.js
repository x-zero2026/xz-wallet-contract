const { Client } = require('pg');

async function main() {
  const client = new Client({
    connectionString: 'postgresql://postgres.rbpsksuuvtzmathnmyxn:iPass4xz2026!@aws-1-ap-south-1.pooler.supabase.com:6543/postgres',
    ssl: {
      rejectUnauthorized: false
    }
  });

  try {
    await client.connect();
    console.log('✅ 数据库连接成功\n');

    // 查询 chilly 的信息
    const result = await client.query(`
      SELECT username, email, did, eth_address 
      FROM users 
      WHERE username = 'chilly' OR email = 'chilly.zhong@gmail.com'
    `);

    if (result.rows.length === 0) {
      console.log('❌ 未找到 chilly 用户');
    } else {
      console.log('========== Chilly 用户信息 ==========');
      result.rows.forEach(user => {
        console.log(`用户名: ${user.username}`);
        console.log(`邮箱: ${user.email}`);
        console.log(`DID: ${user.did}`);
        console.log(`ETH地址: ${user.eth_address || '(未设置)'}`);
        console.log('');
      });
    }

    // 查询所有有 eth_address 的用户
    const allUsers = await client.query(`
      SELECT username, eth_address 
      FROM users 
      WHERE eth_address IS NOT NULL AND eth_address != ''
      ORDER BY username
    `);

    console.log('========== 所有已设置 ETH 地址的用户 ==========');
    allUsers.rows.forEach(user => {
      console.log(`${user.username}: ${user.eth_address}`);
    });

  } catch (error) {
    console.error('❌ 错误:', error.message);
  } finally {
    await client.end();
  }
}

main();
