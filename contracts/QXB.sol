// SPDX-License-Identifier: MIT
pragma solidity ^0.8.30;

/**
 * @title QXB (齐夏币)
 * @dev 一个简单的 ERC20 代币合约
 * @notice 这是一个学习用的代币合约，部署在 Sepolia 测试网上
 * @notice 代币名称：齐夏币，符号：QXB
 */
contract QXB {
    // 合约版本信息
    string public constant VERSION = "1.0.0";  // 合约版本号
    
    // 代币基本信息
    string public name;           // 代币名称
    string public symbol;         // 代币符号
    uint8 public decimals;        // 小数位数
    uint256 public totalSupply;  // 总供应量

    // 地址余额映射
    mapping(address => uint256) public balanceOf;
    
    // 授权映射：owner => spender => amount
    mapping(address => mapping(address => uint256)) public allowance;

    // 权限管理
    address public owner;

    // 每日奖励相关
    mapping(address => uint256) public lastClaimDay;   // 记录每个地址上次领取的日期（天数，从某个起点开始）
    uint256 public constant DAILY_REWARD = 1e18;      // 每日奖励：1 QXB (1 * 10^18)
    uint256 public constant DAY_IN_SECONDS = 86400;    // 一天的秒数

    // 作者简历（Markdown），仅合约所有者可更新
    string private authorResume;

    // 事件定义
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    event DailyRewardClaimed(address indexed user, uint256 amount, uint256 timestamp);
    event ResumeUpdated(address indexed updater);

    /**
     * @dev 构造函数，初始化代币
     * @param _name 代币名称
     * @param _symbol 代币符号
     * @param _decimals 小数位数
     * @param _totalSupply 总供应量（需要考虑小数位数）
     */
    constructor(
        string memory _name,
        string memory _symbol,
        uint8 _decimals,
        uint256 _totalSupply
    ) {
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        totalSupply = _totalSupply;
        owner = msg.sender;  // 记录部署者地址
        
        // 将所有代币分配给部署者（合约创建者）
        balanceOf[msg.sender] = _totalSupply;
        
        // 触发转账事件（从零地址到部署者）
        emit Transfer(address(0), msg.sender, _totalSupply);
    }

    /**
     * @dev 设置作者简历（Markdown），仅合约所有者可调用
     */
    function setResume(string calldata _resume) external {
        require(msg.sender == owner, "Only owner");
        authorResume = _resume;
        emit ResumeUpdated(msg.sender);
    }

    /**
     * @dev 读取作者简历（Markdown）
     */
    function getResume() external view returns (string memory) {
        return authorResume;
    }

    /**
     * @dev 转账代币
     * @param _to 接收地址
     * @param _value 转账数量
     * @return success 是否成功
     */
    function transfer(address _to, uint256 _value) public returns (bool success) {
        // 检查余额是否足够
        require(balanceOf[msg.sender] >= _value, "Insufficient balance");
        require(_value > 0, "Transfer amount must be greater than 0");
        
        // 检查接收地址是否有效（不能是零地址）
        require(_to != address(0), "Invalid recipient address");
        
        // 执行转账
        balanceOf[msg.sender] -= _value;
        balanceOf[_to] += _value;
        
        // 触发转账事件
        emit Transfer(msg.sender, _to, _value);
        
        return true;
    }

    /**
     * @dev 授权其他地址使用你的代币
     * @notice 标准 ERC20 approve 实现
     * @param _spender 被授权的地址
     * @param _value 授权数量
     * @return success 是否成功
     */
    function approve(address _spender, uint256 _value) public returns (bool success) {
        require(_spender != address(0), "Invalid spender address");
        
        // 设置授权额度（标准 ERC20 实现）
        allowance[msg.sender][_spender] = _value;
        
        // 触发授权事件
        emit Approval(msg.sender, _spender, _value);
        
        return true;
    }

    /**
     * @dev 增加授权额度（推荐使用，避免前置攻击风险）
     * @param _spender 被授权的地址
     * @param _addedValue 增加的授权数量
     * @return success 是否成功
     */
    function increaseAllowance(address _spender, uint256 _addedValue) public returns (bool success) {
        require(_spender != address(0), "Invalid spender address");
        
        uint256 currentAllowance = allowance[msg.sender][_spender];
        uint256 newAllowance = currentAllowance + _addedValue;
        
        // 检查溢出（虽然 Solidity 0.8+ 会自动检查，但显式检查更清晰）
        require(newAllowance >= currentAllowance, "Allowance overflow");
        
        allowance[msg.sender][_spender] = newAllowance;
        emit Approval(msg.sender, _spender, newAllowance);
        
        return true;
    }

    /**
     * @dev 减少授权额度
     * @param _spender 被授权的地址
     * @param _subtractedValue 减少的授权数量
     * @return success 是否成功
     */
    function decreaseAllowance(address _spender, uint256 _subtractedValue) public returns (bool success) {
        require(_spender != address(0), "Invalid spender address");
        
        uint256 currentAllowance = allowance[msg.sender][_spender];
        require(currentAllowance >= _subtractedValue, "Decreased allowance below zero");
        
        uint256 newAllowance = currentAllowance - _subtractedValue;
        allowance[msg.sender][_spender] = newAllowance;
        emit Approval(msg.sender, _spender, newAllowance);
        
        return true;
    }

    /**
     * @dev 从授权地址转账代币（代理转账）
     * @param _from 代币来源地址
     * @param _to 接收地址
     * @param _value 转账数量
     * @return success 是否成功
     */
    function transferFrom(address _from, address _to, uint256 _value) public returns (bool success) {
        // 检查授权额度是否足够
        require(allowance[_from][msg.sender] >= _value, "Insufficient allowance");
        
        // 检查来源地址余额是否足够
        require(balanceOf[_from] >= _value, "Insufficient balance");
        require(_value > 0, "Transfer amount must be greater than 0");
        
        // 检查接收地址是否有效
        require(_to != address(0), "Invalid recipient address");
        
        // 减少授权额度
        allowance[_from][msg.sender] -= _value;
        
        // 执行转账
        balanceOf[_from] -= _value;
        balanceOf[_to] += _value;
        
        // 触发转账事件
        emit Transfer(_from, _to, _value);
        
        return true;
    }

    /**
     * @dev 铸造新代币（仅合约所有者可以调用）
     * @notice 修复：使用 owner 变量而不是 tx.origin 进行权限控制
     * @param _to 接收地址
     * @param _amount 铸造数量
     */
    function mint(address _to, uint256 _amount) public {
        // 修复：使用 owner 变量进行权限控制
        require(msg.sender == owner, "Only owner can mint");
        require(_to != address(0), "Invalid recipient address");
        require(_amount > 0, "Amount must be greater than 0");
        
        // 增加总供应量
        totalSupply += _amount;
        
        // 增加接收地址余额
        balanceOf[_to] += _amount;
        
        // 触发转账事件（从零地址到接收地址）
        emit Transfer(address(0), _to, _amount);
    }

    /**
     * @dev 销毁代币
     * @param _amount 销毁数量
     */
    function burn(uint256 _amount) public {
        // 检查余额是否足够
        require(balanceOf[msg.sender] >= _amount, "Insufficient balance");
        require(_amount > 0, "Amount must be greater than 0");
        
        // 减少总供应量
        totalSupply -= _amount;
        
        // 减少发送者余额
        balanceOf[msg.sender] -= _amount;
        
        // 触发转账事件（从发送者到零地址）
        emit Transfer(msg.sender, address(0), _amount);
    }

    /**
     * @dev 内部函数：检查用户是否可以领取每日奖励（基于日期）
     * @param _user 要检查的地址
     * @return canClaim 是否可以领取
     * @return nextClaimDay 下次可以领取的日期（天数，如果当前可以领取，返回0）
     */
    function _checkCanClaim(address _user) internal view returns (bool canClaim, uint256 nextClaimDay) {
        // 计算当前日期（天数）：block.timestamp / DAY_IN_SECONDS
        uint256 currentDay = block.timestamp / DAY_IN_SECONDS;
        uint256 lastClaim = lastClaimDay[_user];
        
        // 如果从未领取过，或上次领取的日期不是今天，可以领取
        if (lastClaim == 0 || lastClaim != currentDay) {
            return (true, 0);
        }
        
        // 今天已经领取过，返回明天可以领取
        return (false, currentDay + 1);
    }

    /**
     * @dev 领取每日奖励
     * @notice 每个地址每天可以领取1枚QXB代币（基于UTC日期）
     * @notice 任何人都可以领取，无需持有代币
     * @notice 使用日期（天数）作为key，不关心具体时间
     * @return success 是否成功
     */
    function claimDailyReward() public returns (bool success) {
        address user = msg.sender;
        
        // 计算当前日期（天数）：block.timestamp / DAY_IN_SECONDS
        uint256 currentDay = block.timestamp / DAY_IN_SECONDS;
        
        // 检查今天是否已经领取过
        require(lastClaimDay[user] != currentDay, "Today's reward already claimed");
        
        // 记录今天已领取（使用当前日期作为key）
        lastClaimDay[user] = currentDay;
        
        // 获取当前时间戳（用于事件）
        uint256 currentTime = block.timestamp;
        
        // 铸造奖励代币
        totalSupply += DAILY_REWARD;
        balanceOf[user] += DAILY_REWARD;
        
        // 触发转账事件（从零地址到用户）
        emit Transfer(address(0), user, DAILY_REWARD);
        
        // 触发每日奖励事件
        emit DailyRewardClaimed(user, DAILY_REWARD, currentTime);
        
        return true;
    }

    /**
     * @dev 查询是否可以领取每日奖励
     * @param _user 查询的地址
     * @return canClaim 是否可以领取
     * @return nextClaimDay 下次可以领取的日期（天数，如果当前可以领取，返回0）
     */
    function canClaimDailyReward(address _user) public view returns (bool canClaim, uint256 nextClaimDay) {
        // 复用内部检查函数
        return _checkCanClaim(_user);
    }

    /**
     * @dev 查询用户上次领取的日期（天数）
     * @param _user 查询的地址
     * @return lastDay 上次领取的日期（天数，0表示从未领取）
     * @return currentDay 当前日期（天数）
     */
    function getClaimDayInfo(address _user) public view returns (uint256 lastDay, uint256 currentDay) {
        return (lastClaimDay[_user], block.timestamp / DAY_IN_SECONDS);
    }
}

